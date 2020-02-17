package helmclient

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

// GetReleaseHistory gets the current installed version of the Helm Release.
// The releaseName is the name of the Helm Release that is set when the Helm
// Chart is installed.
func (c *Client) GetReleaseHistory(ctx context.Context, namespace, releaseName string) (*ReleaseHistory, error) {
	eventName := "get_release_history"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	releaseContent, err := c.getReleaseHistory(ctx, namespace, releaseName)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return nil, microerror.Mask(err)
	}

	return releaseContent, nil
}

func (c *Client) getReleaseHistory(ctx context.Context, namespace, releaseName string) (*ReleaseHistory, error) {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	history := action.NewHistory(cfg)

	// We only care about the most recent history record.
	history.Max = 1

	releases, err := history.Run(releaseName)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if len(releases) == 0 {
		return nil, nil
	}

	return releaseToReleaseHistory(releases), nil
}

func releaseToReleaseHistory(releases []*release.Release) *ReleaseHistory {
	release := releases[0]

	var appVersion, description, version string

	if release.Chart != nil && release.Chart.Metadata != nil {
		appVersion = release.Chart.Metadata.AppVersion
		version = release.Chart.Metadata.Version
	}

	var lastDeployed time.Time

	if release.Info != nil {
		description = release.Info.Description
		lastDeployed = release.Info.LastDeployed.Time
	}

	return &ReleaseHistory{
		AppVersion:   appVersion,
		Description:  description,
		LastDeployed: lastDeployed,
		Name:         release.Name,
		Version:      version,
	}
}
