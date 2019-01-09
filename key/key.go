package key

import (
	hapichart "k8s.io/helm/pkg/proto/hapi/chart"
)

func ChartVersion(chart *hapichart.Chart) string {
	if chart.Metadata == nil {
		return ""
	}

	return chart.Metadata.Version
}
