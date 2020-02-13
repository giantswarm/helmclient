package helmclient

import (
	"context"
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/afero"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/cmd/util"
)

const (
	defaultReleaseNamespace = "giantswarm"
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

	RESTConfig       *rest.Config
	ReleaseNamespace string
}

// Client knows how to talk with a Helm Tiller server.
type Client struct {
	fs         afero.Fs
	helmClient Interface
	httpClient *http.Client
	k8sClient  kubernetes.Interface
	logger     micrologger.Logger

	restConfig       *rest.Config
	releaseNamespace string
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

	if config.RESTConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.RESTConfig must not be empty", config)
	}
	if config.ReleaseNamespace == "" {
		config.ReleaseNamespace = defaultReleaseNamespace
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

		releaseNamespace: config.ReleaseNamespace,
		restConfig:       config.RESTConfig,
	}

	return c, nil
}

// DeleteRelease uninstalls a chart given its release name.
func (c *Client) DeleteRelease(ctx context.Context, releaseName string) error {
	eventName := "delete_release"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.deleteRelease(ctx, releaseName)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) deleteRelease(ctx context.Context, releaseName string) error {
	c.logger.LogCtx(ctx, "level", "debug", "message", "delete release not yet implemented for helm 3")
	return nil
}

// GetReleaseHistory gets the current installed version of the Helm Release.
// The releaseName is the name of the Helm Release that is set when the Helm
// Chart is installed.
func (c *Client) GetReleaseHistory(ctx context.Context, releaseName string) (*ReleaseHistory, error) {
	eventName := "get_release_history"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	releaseContent, err := c.getReleaseHistory(ctx, releaseName)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return nil, microerror.Mask(err)
	}

	return releaseContent, nil
}

func (c *Client) getReleaseHistory(ctx context.Context, releaseName string) (*ReleaseHistory, error) {
	c.logger.LogCtx(ctx, "level", "debug", "message", "get release history not yet implemented for helm 3")
	return nil, nil
}

// InstallReleaseFromTarball installs a chart packaged in the given tarball.
func (c *Client) InstallReleaseFromTarball(ctx context.Context, chartPath string, values map[string]interface{}, options InstallOptions) error {
	eventName := "install_release_from_tarball"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.installReleaseFromTarball(ctx, chartPath, values, options)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) installReleaseFromTarball(ctx context.Context, chartPath string, values map[string]interface{}, options InstallOptions) error {
	c.logger.LogCtx(ctx, "level", "debug", "message", "install release from tarball not yet implemented for helm 3")
	return nil
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

func (c *Client) newActionConfig() (*action.Configuration, error) {
	restClientGetter := newConfigFlags(c.restConfig, c.releaseNamespace)
	kubeClient := &kube.Client{
		Factory: util.NewFactory(restClientGetter),
	}

	// Use secrets driver for release storage.
	s := driver.NewSecrets(c.k8sClient.CoreV1().Secrets(c.releaseNamespace))

	store := storage.Init(s)

	return &action.Configuration{
		RESTClientGetter: restClientGetter,
		Releases:         store,
		KubeClient:       kubeClient,
	}, nil
}

func newConfigFlags(config *rest.Config, namespace string) *genericclioptions.ConfigFlags {
	return &genericclioptions.ConfigFlags{
		Namespace:   &namespace,
		APIServer:   &config.Host,
		CAFile:      &config.CAFile,
		BearerToken: &config.BearerToken,
	}
}
