//go:build k8srequired
// +build k8srequired

package releasetesting

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/helmclient/v4/integration/charttarball"
	"github.com/giantswarm/helmclient/v4/pkg/helmclient"
)

func TestBasic(t *testing.T) {
	ctx := context.Background()

	var passingReleaseName string = "passing-test-chart"

	{
		err := runReleaseTest(ctx, passingReleaseName)
		if err != nil {
			t.Fatalf("expected nil error got %#v", err)
		}
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting release %#q", passingReleaseName))

		err := config.HelmClient.DeleteRelease(ctx, metav1.NamespaceDefault, passingReleaseName)
		if err != nil {
			t.Fatalf("expected nil error got %#v", err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleted release %#q", passingReleaseName))
	}

	var failingReleaseName string = "failing-test-chart"

	{
		err := runReleaseTest(ctx, failingReleaseName)
		if !helmclient.IsTestReleaseFailure(err) {
			t.Fatalf("expected release test failure got %#v", err)
		}
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting release %#q", failingReleaseName))

		err := config.HelmClient.DeleteRelease(ctx, metav1.NamespaceDefault, failingReleaseName)
		if err != nil {
			t.Fatalf("expected nil error got %#v", err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleted release %#q", failingReleaseName))
	}
}

func runReleaseTest(ctx context.Context, releaseName string) error {
	var err error
	var chartPath = ""

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating tarball for %#q", releaseName))

		chartPath, err = charttarball.Create(releaseName)
		if err != nil {
			return microerror.Mask(err)
		}
		defer os.Remove(chartPath)

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created tarball for %#q", releaseName))
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("installing %#q", releaseName))

		installOptions := helmclient.InstallOptions{
			ReleaseName: releaseName,
			Wait:        true,
		}
		err = config.HelmClient.InstallReleaseFromTarball(ctx, chartPath, metav1.NamespaceDefault, map[string]interface{}{}, installOptions)
		if err != nil {
			return microerror.Mask(err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("installed %#q", releaseName))
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("running release tests for %#q", releaseName))

		err = config.HelmClient.RunReleaseTest(ctx, metav1.NamespaceDefault, releaseName)
		if err != nil {
			config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("release test failed for %#q", releaseName))
			return microerror.Mask(err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("release test passed for %#q", releaseName))
	}

	return nil
}
