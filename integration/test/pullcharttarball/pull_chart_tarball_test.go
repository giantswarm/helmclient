// +build k8srequired

package pullcharttarball

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/k8sportforward"
	"github.com/giantswarm/microerror"
	"k8s.io/helm/pkg/helm"
)

// TestPullChartTarball starts a chartmuseum server and pushes a test chart.
// PullChartTarball is used to download the chart and LoadChart is used to
// parse it. Finally the metadata for the chart is checked.
func TestPullChartTarball(t *testing.T) {
	ctx := context.Background()
	var err error

	err = installChartMuseum(ctx)
	if err != nil {
		t.Fatalf("could not install chartmuseum %#v", err)
	}

	var fw *k8sportforward.Forwarder
	{
		c := k8sportforward.ForwarderConfig{
			RestConfig: config.RestConfig,
		}

		fw, err = k8sportforward.NewForwarder(c)
		if err != nil {
			t.Fatalf("could not create forwarder %v", err)
		}
	}

	podName, err := config.Host.GetPodName("default", "app=chartmuseum")
	if err != nil {
		t.Fatalf("could not get chartmuseum pod name %#v", err)
	}
	tunnel, err := fw.ForwardPort("default", podName, 8080)
	if err != nil {
		t.Fatalf("could not create tunnel %#v", err)
	}

	serverAddress := "http://" + tunnel.LocalAddress()
	err = waitForServer(ctx, serverAddress+"/health")
	if err != nil {
		t.Fatalf("server didn't come up on time")
	}

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
	if chart.Version != expectedVersion {
		t.Fatalf("expected chart version %#q got %#q", expectedVersion, chart.Version)
	}
}

func installChartMuseum(ctx context.Context) error {
	chartMuseumRelease := "chartmuseum"
	chartMuseumTarball, err := charttarball.Create("chartmuseum-chart")
	if err != nil {
		return microerror.Mask(err)
	}
	defer os.Remove(chartMuseumTarball)

	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	// We need to pass the ValueOverrides option to make the install process
	// use the default values and prevent errors on nested values.
	err = config.HelmClient.InstallReleaseFromTarball(ctx, chartMuseumTarball, "default", helm.ReleaseName(chartMuseumRelease), helm.ValueOverrides([]byte("{}")))
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func waitForServer(ctx context.Context, url string) error {
	var err error

	o := func() error {
		_, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("could not retrieve %s: %v", url, err)
		}
		return nil
	}

	b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
	n := backoff.NewNotifier(c.logger, ctx)

	err = backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
