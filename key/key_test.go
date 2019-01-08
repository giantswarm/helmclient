package key

import (
	"testing"

	hapichart "k8s.io/helm/pkg/proto/hapi/chart"
)

func Test_ChartVersion(t *testing.T) {
	testCases := []struct {
		name            string
		chart           *hapichart.Chart
		expectedVersion string
	}{
		{
			name: "case 0: basic match",
			chart: &hapichart.Chart{
				Metadata: &hapichart.Metadata{
					Version: "1.0.0",
				},
			},
			expectedVersion: "1.0.0",
		},
		{
			name: "case 1: empty metadata returns empty version",
			chart: &hapichart.Chart{
				Metadata: &hapichart.Metadata{},
			},
			expectedVersion: "",
		},
		{
			name:            "case 2: empty chart returns empty version",
			chart:           &hapichart.Chart{},
			expectedVersion: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ChartVersion(tc.chart)

			if result != tc.expectedVersion {
				t.Fatalf("ChartVersion == %#q, want %#q", result, tc.expectedVersion)
			}
		})
	}
}
