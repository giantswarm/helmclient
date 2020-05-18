package helmclient

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v3/pkg/action"
)

// DeleteRelease uninstalls a chart given its release name.
func (c *Client) DeleteRelease(ctx context.Context, namespace, releaseName string) error {
	eventName := "delete_release"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.deleteRelease(ctx, namespace, releaseName)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) deleteRelease(ctx context.Context, namespace, releaseName string) error {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return microerror.Mask(err)
	}

	uninstall := action.NewUninstall(cfg)

	_, err = uninstall.Run(releaseName)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
