package helmclient

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v4/pkg/action"
)

// Rollback executes a rollback to a previous revision of a Helm release.
func (c *Client) Rollback(ctx context.Context, namespace, releaseName string, revision int, options RollbackOptions) error {
	eventName := "rollback"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer func() {
		eventCounter.WithLabelValues(eventName, releaseName).Inc()
		t.ObserveDuration()
	}()

	err := c.rollback(ctx, namespace, releaseName, revision, options)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) rollback(ctx context.Context, namespace, releaseName string, revision int, options RollbackOptions) error {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return microerror.Mask(err)
	}

	rollback := action.NewRollback(cfg)

	// Configure action with supported rollback options.
	options.configure(rollback, namespace, revision)

	err = rollback.Run(releaseName)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (options RollbackOptions) configure(action *action.Rollback, namespace string, revision int) {
	if options.Timeout == 0 {
		options.Timeout = time.Second * defaultK8sClientTimeout
	}

	action.Force = options.Force
	action.Timeout = options.Timeout
	action.Version = revision
	action.Wait = options.Wait
}
