package helmclient

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v4/pkg/action"
)

// DeleteRelease uninstalls a chart given its release name.
func (c *Client) DeleteRelease(ctx context.Context, namespace, releaseName string, options DeleteOptions) error {
	eventName := "delete_release"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer func() {
		eventCounter.WithLabelValues(eventName, releaseName).Inc()
		t.ObserveDuration()
	}()

	err := c.deleteRelease(ctx, namespace, releaseName, options)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) deleteRelease(ctx context.Context, namespace, releaseName string, options DeleteOptions) error {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return microerror.Mask(err)
	}

	uninstall := action.NewUninstall(cfg)

	// Configure action with supported upgrade options.
	options.configure(uninstall)

	_, err = uninstall.Run(releaseName)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (options DeleteOptions) configure(action *action.Uninstall) {
	if options.Timeout > 0 {
		action.Timeout = options.Timeout
	}
}
