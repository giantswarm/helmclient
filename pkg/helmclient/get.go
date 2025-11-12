package helmclient

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v4/pkg/action"
)

// GetReleaseContent gets the current status of the Helm Release including any
// values provided when the chart was installed. The releaseName is the name
// of the Helm Release that is set when the Helm Chart is installed.
func (c *Client) GetReleaseContent(ctx context.Context, namespace, releaseName string) (*ReleaseContent, error) {
	eventName := "get_release_content"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	releaseContent, err := c.getReleaseContent(ctx, namespace, releaseName)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return nil, microerror.Mask(err)
	}

	return releaseContent, nil
}

func (c *Client) getReleaseContent(ctx context.Context, namespace, releaseName string) (*ReleaseContent, error) {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	get := action.NewGet(cfg)

	res, err := get.Run(releaseName)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return releaseToReleaseContent(res), nil
}
