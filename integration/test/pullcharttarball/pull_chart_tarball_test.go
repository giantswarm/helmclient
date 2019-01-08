// +build k8srequired

package pullcharttarball

import (
	"context"
	"os"
	"testing"

	"github.com/giantswarm/helmclient/integration/charttarball"
	"github.com/giantswarm/helmclient/key"
)

func TestPullChartTarball(t *testing.T) {
	ctx := context.Background()
	var err error

	const releaseName = "test"

	tarballPath, err := charttarball.Create("test-chart")
	if err != nil {
		t.Fatalf("could not create chart archive %#v", err)
	}
	defer os.Remove(tarballPath)

	chart, err := config.HelmClient.LoadChart(ctx, tarballPath)
	if err != nil {
		t.Fatalf("could not load chart %#v", err)
	}

	expectedVersion := "1.2.3"
	if key.ChartVersion(chart) != expectedVersion {
		t.Fatalf("expected chart version %#q got %#q", expectedVersion, key.ChartVersion(chart))
	}
}
