package helmclient

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// UpdateReleaseFromTarball updates the given release using the chart packaged
// in the tarball.
func (c *Client) UpdateReleaseFromTarball(ctx context.Context, chartPath, namespace, releaseName string, values map[string]interface{}, options UpdateOptions) error {
	eventName := "update_release_from_tarball"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.updateReleaseFromTarball(ctx, chartPath, namespace, releaseName, values, options)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) updateReleaseFromTarball(ctx context.Context, chartPath, namespace, releaseName string, values map[string]interface{}, options UpdateOptions) error {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return microerror.Mask(err)
	}

	upgrade := action.NewUpgrade(cfg)

	// Load the chart from the given path. This also ensures that all chart
	// dependencies are present.
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return microerror.Mask(err)
	}

	// Configure action with supported upgrade options.
	options.configure(upgrade, namespace)

	_, err = upgrade.Run(releaseName, chartRequested, values)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (options UpdateOptions) configure(action *action.Upgrade, namespace string) {
	action.Force = options.Force
	action.Namespace = namespace
	action.Wait = options.Wait
}
