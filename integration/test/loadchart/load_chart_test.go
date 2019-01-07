// +build k8srequired

package loadchart

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/helmclient/integration/charttarball"
)

func TestLoadChart(t *testing.T) {
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

	config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("chart metadata %#v", chart.Metadata))
}
