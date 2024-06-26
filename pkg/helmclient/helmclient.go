package helmclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	kubeconfig "github.com/giantswarm/kubeconfig/v4"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/afero"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"oras.land/oras-go/pkg/content"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

// Config represents the configuration used to create a helm client.
type Config struct {
	Fs afero.Fs
	// HelmClient sets a helm client used for all operations of the initiated
	// client. If this is nil, a new helm client will be created. Setting the
	// helm client here manually might only be sufficient for testing or
	// whenever you know what you do.
	HelmClient Interface
	K8sClient  kubernetes.Interface
	Logger     micrologger.Logger
	// RegistryOptions can be used to allow pulling from private chart
	// registries or accessing them insecurely. If this is nil, empty options
	// will be created. This will result in using default helm config found on
	// the machine (~/.config/helm). We assume the registries we pull from are
	// public.
	RegistryOptions *content.RegistryOptions
	RestClient      rest.Interface
	RestConfig      *rest.Config
	RestMapper      meta.RESTMapper

	HTTPClientTimeout time.Duration
}

// Client knows how to talk with Helm.
type Client struct {
	fs              afero.Fs
	helmClient      Interface
	httpClient      *http.Client
	k8sClient       kubernetes.Interface
	logger          micrologger.Logger
	registryOptions content.RegistryOptions
	restClient      rest.Interface
	restConfig      *rest.Config
	restMapper      meta.RESTMapper
}

// debugLogFunc allows us to pass log messages from helm to micrologger.
type debugLogFunc func(string, ...interface{})

// restClientGetter gets a REST client for use by the Helm kube client.
type restClientGetter struct {
	discoveryClient     discovery.CachedDiscoveryInterface
	rawKubeConfigLoader clientcmd.ClientConfig
	restConfig          *rest.Config
	restMapper          meta.RESTMapper
}

// New creates a new configured Helm client.
func New(config Config) (*Client, error) {
	if config.Fs == nil {
		config.Fs = afero.NewOsFs()
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.RegistryOptions == nil {
		config.RegistryOptions = &content.RegistryOptions{
			Configs:   []string{},
			Username:  "",
			Password:  "",
			Insecure:  false,
			PlainHTTP: false,
		}
	}
	if config.RestClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.RestClient must not be empty", config)
	}
	if config.RestConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.RestConfig must not be empty", config)
	}

	rmHttpClient, err := rest.HTTPClientFor(rest.CopyConfig(config.RestConfig))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if config.RestMapper == nil {
		restMapper, err := apiutil.NewDynamicRESTMapper(rest.CopyConfig(config.RestConfig), rmHttpClient)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		config.RestMapper = restMapper
	}

	if config.HTTPClientTimeout == 0 {
		config.HTTPClientTimeout = defaultHTTPClientTimeout
	}

	// Set client timeout to prevent leakages.
	httpClient := &http.Client{
		Timeout: time.Second * time.Duration(config.HTTPClientTimeout),
	}

	c := &Client{
		fs:              config.Fs,
		helmClient:      config.HelmClient,
		httpClient:      httpClient,
		k8sClient:       config.K8sClient,
		logger:          config.Logger,
		registryOptions: *config.RegistryOptions,
		restClient:      config.RestClient,
		restConfig:      config.RestConfig,
		restMapper:      config.RestMapper,
	}

	return c, nil
}

// debugLogFunc allows us to pass micrologger to components that expect a
// klog.Infof function. We downgrade the messages from info to debug to match
// our usual approach.
func (c *Client) debugLogFunc(ctx context.Context) debugLogFunc {
	return func(format string, args ...interface{}) {
		message := fmt.Sprintf(format, args...)
		c.logger.LogCtx(ctx, "level", "debug", "message", message)
	}
}

// newActionConfig creates a config for the Helm action package.
func (c *Client) newActionConfig(ctx context.Context, namespace string) (*action.Configuration, error) {
	restClient, err := c.newRESTClientGetter(ctx, namespace)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Create a Helm kube client.
	kubeClient := kube.New(restClient)

	// Use secrets driver for release storage.
	s := driver.NewSecrets(c.k8sClient.CoreV1().Secrets(namespace))
	store := storage.Init(s)

	return &action.Configuration{
		Log:              c.debugLogFunc(ctx),
		KubeClient:       kubeClient,
		Releases:         store,
		RESTClientGetter: restClient,
	}, nil
}

func (c *Client) newRESTClientGetter(ctx context.Context, namespace string) (*restClientGetter, error) {
	if namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "namespace must not be empty")
	}

	// Create a discovery client using the in memory cache.
	discoveryClient := discovery.NewDiscoveryClient(c.restClient)
	cachedDiscoveryClient := memory.NewMemCacheClient(discoveryClient)

	// Convert REST config back to a kubeconfig for the raw kubeconfig loader.
	bytes, err := kubeconfig.NewKubeConfigForRESTConfig(ctx, c.restConfig, "helmclient", namespace)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	rawKubeConfigLoader, err := clientcmd.NewClientConfigFromBytes(bytes)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &restClientGetter{
		discoveryClient:     cachedDiscoveryClient,
		rawKubeConfigLoader: rawKubeConfigLoader,
		restConfig:          c.restConfig,
		restMapper:          c.restMapper,
	}, nil
}

func (r *restClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return r.discoveryClient, nil
}

func (r *restClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return r.rawKubeConfigLoader
}

func (r *restClientGetter) ToRESTConfig() (*rest.Config, error) {
	return r.restConfig, nil
}

func (r *restClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	return r.restMapper, nil
}

func releaseToReleaseContent(res *release.Release) *ReleaseContent {
	release := &ReleaseContent{
		Name:     res.Name,
		Revision: res.Version,
		Status:   res.Info.Status.String(),
		Values:   res.Config,
	}

	if res.Chart != nil && res.Chart.Metadata != nil {
		release.AppVersion = res.Chart.Metadata.AppVersion
		release.Version = res.Chart.Metadata.Version
	}

	if res.Info != nil {
		release.Description = res.Info.Description
		release.LastDeployed = res.Info.LastDeployed.Time
	}

	return release
}
