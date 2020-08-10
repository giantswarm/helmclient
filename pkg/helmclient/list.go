package helmclient

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v3/pkg/action"
)

// ListReleaseContents gets the current status of all Helm Releases.
func (c *Client) ListReleaseContents(ctx context.Context, namespace string) ([]*ReleaseContent, error) {
	eventName := "list_release_contents"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	releaseContent, err := c.listReleaseContents(ctx, namespace)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return nil, microerror.Mask(err)
	}

	return releaseContent, nil
}

func (c *Client) listReleaseContents(ctx context.Context, namespace string) ([]*ReleaseContent, error) {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	list := action.NewList(cfg)

	res, err := list.Run()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var releases = []*ReleaseContent{}

	for _, rel := range res {
		releases = append(releases, releaseToReleaseContent(rel))
	}

	return releases, nil
}
