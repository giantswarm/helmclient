package helmclient

import (
	"context"
	"net/http"
	"time"

	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/kubeconfig"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/afero"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	K8sClient  k8sclient.Interface
	Logger     micrologger.Logger
}

// Client knows how to talk with Helm.
type Client struct {
	fs         afero.Fs
	helmClient Interface
	httpClient *http.Client
	k8sClient  k8sclient.Interface
	logger     micrologger.Logger
}

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
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}

	// Set client timeout to prevent leakages.
	httpClient := &http.Client{
		Timeout: time.Second * httpClientTimeout,
	}

	c := &Client{
		fs:         config.Fs,
		helmClient: config.HelmClient,
		httpClient: httpClient,
		k8sClient:  config.K8sClient,
		logger:     config.Logger,
	}

	return c, nil
}

// ListReleaseContents gets the current status of all Helm Releases.
func (c *Client) ListReleaseContents(ctx context.Context) ([]*ReleaseContent, error) {
	eventName := "list_release_contents"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	releaseContent, err := c.listReleaseContents(ctx)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return nil, microerror.Mask(err)
	}

	return releaseContent, nil
}

func (c *Client) listReleaseContents(ctx context.Context) ([]*ReleaseContent, error) {
	c.logger.LogCtx(ctx, "level", "debug", "message", "list release contents not yet implemented for helm 3")
	return nil, nil
}

// LoadChart loads a Helm Chart and returns relevant parts of its structure.
func (c *Client) LoadChart(ctx context.Context, chartPath string) (Chart, error) {
	eventName := "load_chart"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	chart, err := c.loadChart(ctx, chartPath)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return Chart{}, microerror.Mask(err)
	}

	return chart, nil
}

func (c *Client) loadChart(ctx context.Context, chartPath string) (Chart, error) {
	c.logger.LogCtx(ctx, "level", "debug", "message", "load chart not yet implemented for helm 3")
	return Chart{}, nil
}

// RunReleaseTest runs the tests for a Helm Release. The releaseName is the
// name of the Helm Release that is set when the Helm Chart is installed. This
// is the same action as running the helm test command.
func (c *Client) RunReleaseTest(ctx context.Context, releaseName string, options ReleaseTestOptions) error {
	eventName := "run_release_test"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.runReleaseTest(ctx, releaseName, options)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) runReleaseTest(ctx context.Context, releaseName string, options ReleaseTestOptions) error {
	c.logger.LogCtx(ctx, "level", "debug", "message", "run release test not yet implemented for helm 3")
	return nil
}

// UpdateReleaseFromTarball updates the given release using the chart packaged
// in the tarball.
func (c *Client) UpdateReleaseFromTarball(ctx context.Context, releaseName, chartPath string, values map[string]interface{}, options UpdateOptions) error {
	eventName := "update_release_from_tarball"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.updateReleaseFromTarball(ctx, releaseName, chartPath, values, options)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) updateReleaseFromTarball(ctx context.Context, releaseName, chartPath string, values map[string]interface{}, options UpdateOptions) error {
	c.logger.LogCtx(ctx, "level", "debug", "message", "update release from tarball not yet implemented for helm 3")
	return nil
}

// newActionConfig creates a config for the Helm action package.
func (c *Client) newActionConfig(ctx context.Context, namespace string) (*action.Configuration, error) {
	restClient, err := newRESTClientGetter(ctx, c.k8sClient, namespace)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Create a Helm kube client.
	kubeClient := kube.New(restClient)

	// Use secrets driver for release storage.
	s := driver.NewSecrets(c.k8sClient.K8sClient().CoreV1().Secrets(namespace))
	store := storage.Init(s)

	return &action.Configuration{
		KubeClient:       kubeClient,
		Releases:         store,
		RESTClientGetter: restClient,
	}, nil
}

func newRESTClientGetter(ctx context.Context, k8sClient k8sclient.Interface, namespace string) (*restClientGetter, error) {
	if k8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "k8sClient must not be empty")
	}

	if namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "namespace must not be empty")
	}

	// Create a discovery client using the in memory cache.
	discoveryClient := discovery.NewDiscoveryClient(k8sClient.RESTClient())
	cachedDiscoveryClient := memory.NewMemCacheClient(discoveryClient)

	restMapper, err := apiutil.NewDynamicRESTMapper(rest.CopyConfig(k8sClient.RESTConfig()))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Convert REST config back to a kubeconfig for the raw kubeconfig loader.
	bytes, err := kubeconfig.NewKubeConfigForRESTConfig(ctx, k8sClient.RESTConfig(), "helmclient", namespace)
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
		restConfig:          k8sClient.RESTConfig(),
		restMapper:          restMapper,
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
