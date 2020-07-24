package helmclient

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// InstallReleaseFromTarball installs a chart packaged in the given tarball.
func (c *Client) InstallReleaseFromTarball(ctx context.Context, chartPath, namespace string, values map[string]interface{}, options InstallOptions) error {
	eventName := "install_release_from_tarball"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.installReleaseFromTarball(ctx, chartPath, namespace, values, options)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) installReleaseFromTarball(ctx context.Context, chartPath, namespace string, values map[string]interface{}, options InstallOptions) error {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return microerror.Mask(err)
	}

	install := action.NewInstall(cfg)

	// Load the chart from the given path. This also ensures that all chart
	// dependencies are present.
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return microerror.Mask(err)
	}

	// Configure action with supported install options.
	options.configure(install, namespace)

	_, err = install.Run(chartRequested, values)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (options InstallOptions) configure(action *action.Install, namespace string) {
	if options.Timeout == 0 {
		options.Timeout = time.Second * defaultK8sClientTimeout
	}

	action.Namespace = namespace
	action.ReleaseName = options.ReleaseName
	action.Timeout = options.Timeout
	action.Wait = options.Wait
}
