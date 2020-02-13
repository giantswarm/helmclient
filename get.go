package helmclient

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v3/pkg/action"
)

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
	cfg, err := c.newActionConfig()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	get := action.NewGet(cfg)

	_, err = get.Run(releaseName)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return nil, nil
}
