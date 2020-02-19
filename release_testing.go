package helmclient

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v3/pkg/action"
)

// RunReleaseTest runs the tests for a Helm Release. The releaseName is the
// name of the Helm Release that is set when the Helm Chart is installed. This
// is the same action as running the helm test command.
func (c *Client) RunReleaseTest(ctx context.Context, namespace, releaseName string) error {
	eventName := "run_release_test"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	err := c.runReleaseTest(ctx, namespace, releaseName)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) runReleaseTest(ctx context.Context, namespace, releaseName string) error {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return microerror.Mask(err)
	}

	releaseTesting := action.NewReleaseTesting(cfg)
	releaseTesting.Namespace = namespace

	res, err := releaseTesting.Run(releaseName)
	if err != nil {
		return microerror.Mask(err)
	}

	c.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("RELEASE %#v", res))

	return nil
}
