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
func (c *Client) GetReleaseHistory(ctx context.Context, namespace, releaseName string) ([]ReleaseHistory, error) {
	eventName := "get_release_history"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	releaseHistory, err := c.getReleaseHistory(ctx, namespace, releaseName)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return nil, microerror.Mask(err)
	}

	return releaseHistory, nil
}

func (c *Client) getReleaseHistory(ctx context.Context, namespace, releaseName string) ([]ReleaseHistory, error) {
	cfg, err := c.newActionConfig(ctx, namespace)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	history := action.NewHistory(cfg)

	// We get the 10 most recent Helm releases.
	history.Max = maxHistory

	releases, err := history.Run(releaseName)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if len(releases) == 0 {
		return nil, nil
	}

	return releasesToReleaseHistory(releases), nil
}

func releasesToReleaseHistory(releases []*release.Release) []ReleaseHistory {
	var history []ReleaseHistory

	for _, release := range releases {
		var appVersion, description, version string
		var lastDeployed time.Time

		if release.Chart != nil && release.Chart.Metadata != nil {
			appVersion = release.Chart.Metadata.AppVersion
			version = release.Chart.Metadata.Version
		}

		if release.Info != nil {
			description = release.Info.Description
			lastDeployed = release.Info.LastDeployed.Time
		}

		hist := ReleaseHistory{
			AppVersion:   appVersion,
			Description:  description,
			LastDeployed: lastDeployed,
			Name:         release.Name,
			Version:      version,
		}

		history = append(history, hist)
	}

	return history
}
