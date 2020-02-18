package helmclient

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// LoadChart loads a Helm Chart and returns relevant parts of its structure.
func (c *Client) LoadChart(ctx context.Context, chartPath string) (Chart, error) {
	eventName := "load_chart"

	t := prometheus.NewTimer(histogram.WithLabelValues(eventName))
	defer t.ObserveDuration()

	chart, err := c.loadChart(ctx, chartPath)
	if err != nil {
		errorGauge.WithLabelValues(eventName).Inc()
		return Chart{}, microerror.Mask(err)
	}

	return chart, nil
}

func (c *Client) loadChart(ctx context.Context, chartPath string) (Chart, error) {
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return Chart{}, microerror.Mask(err)
	}

	return newChart(chartRequested)
}

func newChart(helmChart *chart.Chart) (Chart, error) {
	if helmChart == nil || helmChart.Metadata == nil {
		return Chart{}, microerror.Maskf(executionFailedError, "expected non nil argument but got %#v", helmChart)
	}

	chart := Chart{
		Version: helmChart.Metadata.Version,
	}

	return chart, nil
}
