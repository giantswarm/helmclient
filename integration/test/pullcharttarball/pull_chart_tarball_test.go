// +build k8srequired

package pullcharttarball

import (
	"context"
	"fmt"
	"os"
	"testing"

	"k8s.io/helm/pkg/helm"

	"github.com/giantswarm/helmclient/integration/charttarball"
	"github.com/giantswarm/helmclient/key"
)

// TestPullChartTarball starts a chartmuseum server and pushes a test chart.
// PullChartTarball is used to download the chart and LoadChart is used to
// parse it. Finally the metadata for the chart is checked.
func TestPullChartTarball(t *testing.T) {
	ctx := context.Background()
	var err error

	chartMuseumRelease := "chartmuseum"
	chartMuseumTarball, err := charttarball.Create("chartmuseum-chart")
	if err != nil {
		t.Fatalf("could not create chartmuseum archive %#v", err)
	}
	defer os.Remove(chartMuseumTarball)

	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		t.Fatalf("could not install Tiller %#v", err)
	}

	// We need to pass the ValueOverrides option to make the install process
	// use the default values and prevent errors on nested values.
	err = config.HelmClient.InstallReleaseFromTarball(ctx, chartMuseumTarball, "default", helm.ReleaseName(chartMuseumRelease), helm.ValueOverrides([]byte("{}")))
	if err != nil {
		t.Fatalf("failed to install release %#q %#v", chartMuseumRelease, err)
	}

	tarballPath, err := charttarball.Create("test-chart")
	if err != nil {
		t.Fatalf("could not create chart archive %#v", err)
	}
	defer os.Remove(tarballPath)

	config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("chart tarball %#q", tarballPath))

	chart, err := config.HelmClient.LoadChart(ctx, tarballPath)
	if err != nil {
		t.Fatalf("could not load chart %#v", err)
	}

	expectedVersion := "1.2.3"
	if key.ChartVersion(chart) != expectedVersion {
		t.Fatalf("expected chart version %#q got %#q", expectedVersion, key.ChartVersion(chart))
	}
}
