package helmclient

import (
	"context"
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/afero"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/chartutil"
	helmclient "k8s.io/helm/pkg/helm"
	hapichart "k8s.io/helm/pkg/proto/hapi/chart"
)

// Config represents the configuration used to create a helm client.
type Config struct {
	Fs afero.Fs
	// HelmClient sets a helm client used for all operations of the initiated
	// client. If this is nil, a new helm client will be created. Setting the
	// helm client here manually might only be sufficient for testing or
	// whenever you know what you do.
	HelmClient helmclient.Interface
	K8sClient  kubernetes.Interface
	Logger     micrologger.Logger

	RestConfig *rest.Config
}

// Client knows how to talk with a Helm Tiller server.
type Client struct {
	fs         afero.Fs
	helmClient helmclient.Interface
	httpClient *http.Client
	k8sClient  kubernetes.Interface
	logger     micrologger.Logger

	restConfig *rest.Config
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

	if config.RestConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.RestConfig must not be empty", config)
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

		restConfig: config.RestConfig,
	}

	return c, nil
}

// DeleteRelease uninstalls a chart given its release name.
func (c *Client) DeleteRelease(ctx context.Context, releaseName string, options ...helmclient.DeleteOption) error {
	eventName := "delete_release"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.deleteRelease(ctx, releaseName, options...)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) deleteRelease(ctx context.Context, releaseName string, options ...helmclient.DeleteOption) error {
	c.logger.LogCtx(ctx, "level", "debug", "message", "delete release not yet implemented for helm 3")
	return nil
}

// GetReleaseContent gets the current status of the Helm Release including any
// values provided when the chart was installed. The releaseName is the name
// of the Helm Release that is set when the Helm Chart is installed.
func (c *Client) GetReleaseContent(ctx context.Context, releaseName string) (*ReleaseContent, error) {
	eventName := "get_release_content"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	releaseContent, err := c.getReleaseContent(ctx, releaseName)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return nil, microerror.Mask(err)
	}

	return releaseContent, nil
}

func (c *Client) getReleaseContent(ctx context.Context, releaseName string) (*ReleaseContent, error) {
	c.logger.LogCtx(ctx, "level", "debug", "message", "get release content not yet implemented for helm 3")
	return nil, nil
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
func (c *Client) InstallReleaseFromTarball(ctx context.Context, path, ns string, options ...helmclient.InstallOption) error {
	eventName := "install_release_from_tarball"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.installReleaseFromTarball(ctx, path, ns, options...)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) installReleaseFromTarball(ctx context.Context, path, ns string, options ...helmclient.InstallOption) error {
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
	helmChart, err := chartutil.Load(chartPath)
	if err != nil {
		return Chart{}, microerror.Mask(err)
	}

	chart, err := newChart(helmChart)
	if err != nil {
		return Chart{}, microerror.Mask(err)
	}

	return chart, nil
}

// RunReleaseTest runs the tests for a Helm Release. The releaseName is the
// name of the Helm Release that is set when the Helm Chart is installed. This
// is the same action as running the helm test command.
func (c *Client) RunReleaseTest(ctx context.Context, releaseName string, options ...helmclient.ReleaseTestOption) error {
	eventName := "run_release_test"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.runReleaseTest(ctx, releaseName, options...)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) runReleaseTest(ctx context.Context, releaseName string, options ...helmclient.ReleaseTestOption) error {
	c.logger.LogCtx(ctx, "level", "debug", "message", "run release test not yet implemented for helm 3")
	return nil
}

// UpdateReleaseFromTarball updates the given release using the chart packaged
// in the tarball.
func (c *Client) UpdateReleaseFromTarball(ctx context.Context, releaseName, path string, options ...helmclient.UpdateOption) error {
	eventName := "update_release_from_tarball"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.updateReleaseFromTarball(ctx, releaseName, path, options...)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) updateReleaseFromTarball(ctx context.Context, releaseName, path string, options ...helmclient.UpdateOption) error {
	c.logger.LogCtx(ctx, "level", "debug", "message", "update release from tarball not yet implemented for helm 3")
	return nil
}

func newChart(helmChart *hapichart.Chart) (Chart, error) {
	if helmChart == nil || helmChart.Metadata == nil {
		return Chart{}, microerror.Maskf(executionFailedError, "expected non nil argument but got %#v", helmChart)
	}

	chart := Chart{
		Version: helmChart.Metadata.Version,
	}

	return chart, nil
}
