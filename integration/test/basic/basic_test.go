// +build k8srequired

package basic

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/helmclient/integration/charttarball"
)

func TestBasic(t *testing.T) {
	ctx := context.Background()

	var err error

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", "checking release not found")

		releaseContent, err := config.HelmClient.GetReleaseContent(ctx, metav1.NamespaceDefault, "no-release-exists")
		if err != nil && !helmclient.IsReleaseNotFound(err) {
			t.Fatalf("expected release not found error got %v", err)
		}
		if releaseContent != nil {
			t.Fatalf("expected nil release got %v", err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", "checked release not found")
	}

	var chartPath string

	{
		chartPath, err = charttarball.Create("test-chart")
		if err != nil {
			t.Fatalf("could not create chart archive %#v", err)
		}
		defer os.Remove(chartPath)
	}

	var releaseName string = "test-chart"

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("installing %#q", releaseName))

		installOptions := helmclient.InstallOptions{
			ReleaseName: releaseName,
			Wait:        true,
		}
		err = config.HelmClient.InstallReleaseFromTarball(ctx, chartPath, metav1.NamespaceDefault, map[string]interface{}{}, installOptions)
		if err != nil {
			t.Fatalf("could not install chart %v", err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("installed %#q", releaseName))
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("getting release content for %#q", releaseName))

		releaseContent, err := config.HelmClient.GetReleaseContent(ctx, metav1.NamespaceDefault, releaseName)
		if err != nil {
			t.Fatalf("expected nil error got %v", err)
		}

		expectedContent := &helmclient.ReleaseContent{
			Name:   releaseName,
			Status: "deployed",
		}
		if !cmp.Equal(releaseContent, expectedContent) {
			t.Fatalf("want matching ReleaseContent \n %s", cmp.Diff(releaseContent, expectedContent))
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("got release content for %#q", releaseName))
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("getting release history for %#q", releaseName))

		releaseHistory, err := config.HelmClient.GetReleaseHistory(ctx, metav1.NamespaceDefault, releaseName)
		if err != nil {
			t.Fatalf("expected nil error got %v", err)
		}

		if releaseHistory.LastDeployed.IsZero() {
			t.Fatalf("expected non zero last deployed got %v", releaseHistory.LastDeployed)
		}
		// Reset to zero for comparison.
		releaseHistory.LastDeployed = time.Time{}

		expectedHistory := &helmclient.ReleaseHistory{
			AppVersion:  "1.2.3",
			Description: "Install complete",
			Name:        releaseName,
			Version:     "3.2.1",
		}
		if !cmp.Equal(releaseHistory, expectedHistory) {
			t.Fatalf("want matching ReleaseHistory \n %s", cmp.Diff(releaseHistory, expectedHistory))
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("got release history for %#q", releaseName))
	}

	values := map[string]interface{}{
		"test": "value",
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updating %#q", releaseName))

		updateOptions := helmclient.UpdateOptions{
			Wait: true,
		}
		err = config.HelmClient.UpdateReleaseFromTarball(ctx, chartPath, metav1.NamespaceDefault, releaseName, values, updateOptions)
		if err != nil {
			t.Fatalf("could not update chart %v", err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updated %#q", releaseName))
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("getting release content for %#q", releaseName))

		releaseContent, err := config.HelmClient.GetReleaseContent(ctx, metav1.NamespaceDefault, releaseName)
		if err != nil {
			t.Fatalf("expected nil error got %v", err)
		}

		expectedContent := &helmclient.ReleaseContent{
			Name:   releaseName,
			Status: "deployed",
			Values: values,
		}
		if !cmp.Equal(releaseContent, expectedContent) {
			t.Fatalf("want matching ReleaseContent \n %s", cmp.Diff(releaseContent, expectedContent))
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("got release content for %#q", releaseName))
	}

	/*
		{
			config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting release %#q", releaseName))

			err := config.HelmClient.DeleteRelease(ctx, metav1.NamespaceDefault, releaseName)
			if err != nil {
				t.Fatalf("expected nil error got %v", err)
			}

			config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleted release %#q", releaseName))
		}
	*/
}
