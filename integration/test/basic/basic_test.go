//go:build k8srequired
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

	"github.com/giantswarm/helmclient/v4/pkg/helmclient"
)

func TestBasic(t *testing.T) {
	ctx := context.Background()

	var err error
	var releaseName string = "test-chart"

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
		tarballURL := "https://giantswarm.github.io/default-catalog/test-app-1.0.0.tgz"
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("pulling tarball %#q", tarballURL))

		chartPath, err = config.HelmClient.PullChartTarball(ctx, tarballURL)
		if err != nil {
			t.Fatalf("could not pull tarball %#v", err)
		}
		defer os.Remove(chartPath)

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("pulled tarball %#q", tarballURL))
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("loading chart tarball %#q", chartPath))

		chart, err := config.HelmClient.LoadChart(ctx, chartPath)
		if err != nil {
			t.Fatalf("could not load chart %v", err)
		}

		expectedChart := helmclient.Chart{
			Version: "1.0.0",
			Annotations: map[string]string{
				"application.giantswarm.io/metadata":      "https://giantswarm.github.io/default-catalog/test-app-1.0.0.tgz-meta/main.yaml",
				"application.giantswarm.io/readme":        "https://giantswarm.github.io/default-catalog/test-app-1.0.0.tgz-meta/README.md",
				"application.giantswarm.io/team":          "honeybadger",
				"application.giantswarm.io/values-schema": "https://giantswarm.github.io/default-catalog/test-app-1.0.0.tgz-meta/values.schema.json",
			},
		}
		if !cmp.Equal(chart, expectedChart) {
			t.Fatalf("want matching Chart \n %s", cmp.Diff(chart, expectedChart))
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("loaded chart tarball %#q", chartPath))
	}

	values := map[string]interface{}{
		"my": "value",
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("installing %#q", releaseName))

		installOptions := helmclient.InstallOptions{
			ReleaseName: releaseName,
			Wait:        true,
		}
		err = config.HelmClient.InstallReleaseFromTarball(ctx, chartPath, metav1.NamespaceDefault, values, installOptions)
		if err != nil {
			t.Fatalf("could not install chart %v", err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("installed %#q", releaseName))
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", "listing releases")

		releases, err := config.HelmClient.ListReleaseContents(ctx, metav1.NamespaceDefault)
		if err != nil {
			t.Fatalf("could not list releases %v", err)
		}
		if len(releases) != 1 {
			t.Fatalf("expected 1 Releases got \n %d", len(releases))
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", "listed releases")
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("getting release content for %#q", releaseName))

		releaseContent, err := config.HelmClient.GetReleaseContent(ctx, metav1.NamespaceDefault, releaseName)
		if err != nil {
			t.Fatalf("expected nil error got %v", err)
		}

		expectedContent := &helmclient.ReleaseContent{
			AppVersion:  "v2.13.0",
			Description: "Install complete",
			Name:        releaseName,
			Revision:    1,
			Status:      helmclient.StatusDeployed,
			Values:      values,
			Version:     "1.0.0",
		}

		if releaseContent.LastDeployed.IsZero() {
			t.Fatalf("expected non zero last deployed got %v", releaseContent.LastDeployed)
		}
		// Reset to zero for comparison.
		releaseContent.LastDeployed = time.Time{}

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

		if len(releaseHistory) != 1 {
			t.Fatalf("expected 1 history record got %d", len(releaseHistory))
		}
		if releaseHistory[0].LastDeployed.IsZero() {
			t.Fatalf("expected non zero last deployed got %v", releaseHistory[0].LastDeployed)
		}
		// Reset to zero for comparison.
		releaseHistory[0].LastDeployed = time.Time{}

		expectedHistory := []helmclient.ReleaseHistory{
			{
				AppVersion:  "v2.13.0",
				Description: "Install complete",
				Name:        releaseName,
				Revision:    1,
				Status:      helmclient.StatusDeployed,
				Version:     "1.0.0",
			},
		}
		if !cmp.Equal(releaseHistory, expectedHistory) {
			t.Fatalf("want matching ReleaseHistory \n %s", cmp.Diff(releaseHistory, expectedHistory))
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("got release history for %#q", releaseName))
	}

	var updatedChartPath string

	{
		tarballURL := "https://giantswarm.github.io/default-catalog/test-app-1.0.0.tgz"
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("pulling tarball %#q", tarballURL))

		updatedChartPath, err = config.HelmClient.PullChartTarball(ctx, tarballURL)
		if err != nil {
			t.Fatalf("could not pull tarball %#v", err)
		}
		defer os.Remove(chartPath)

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("pulled tarball %#q", tarballURL))
	}

	updatedValues := map[string]interface{}{
		"another": "value",
		"my":      "value",
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updating %#q", releaseName))

		updateOptions := helmclient.UpdateOptions{
			Wait: true,
		}
		err = config.HelmClient.UpdateReleaseFromTarball(ctx, updatedChartPath, metav1.NamespaceDefault, releaseName, updatedValues, updateOptions)
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
			AppVersion:  "v2.13.0",
			Description: "Upgrade complete",
			Name:        releaseName,
			Revision:    2,
			Status:      helmclient.StatusDeployed,
			Values:      updatedValues,
			Version:     "1.0.0",
		}

		if releaseContent.LastDeployed.IsZero() {
			t.Fatalf("expected non zero last deployed got %v", releaseContent.LastDeployed)
		}
		// Reset to zero for comparison.
		releaseContent.LastDeployed = time.Time{}

		if !cmp.Equal(releaseContent, expectedContent) {
			t.Fatalf("want matching ReleaseContent \n %s", cmp.Diff(releaseContent, expectedContent))
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("got release content for %#q", releaseName))
	}

	{
		revision := 1
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("rolling back %#q to revision %d", releaseName, revision))

		rollbackOptions := helmclient.RollbackOptions{
			Wait: true,
		}
		err = config.HelmClient.Rollback(ctx, metav1.NamespaceDefault, releaseName, revision, rollbackOptions)
		if err != nil {
			t.Fatalf("could not rollback %v", err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("rolled back %#q", releaseName))
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("getting release content for %#q", releaseName))

		releaseContent, err := config.HelmClient.GetReleaseContent(ctx, metav1.NamespaceDefault, releaseName)
		if err != nil {
			t.Fatalf("expected nil error got %v", err)
		}

		expectedContent := &helmclient.ReleaseContent{
			AppVersion:  "v2.13.0",
			Description: "Rollback to 1",
			Name:        releaseName,
			Revision:    3,
			Status:      helmclient.StatusDeployed,
			Values:      values,
			Version:     "1.0.0",
		}

		if releaseContent.LastDeployed.IsZero() {
			t.Fatalf("expected non zero last deployed got %v", releaseContent.LastDeployed)
		}
		// Reset to zero for comparison.
		releaseContent.LastDeployed = time.Time{}

		if !cmp.Equal(releaseContent, expectedContent) {
			t.Fatalf("want matching ReleaseContent \n %s", cmp.Diff(releaseContent, expectedContent))
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("got release content for %#q", releaseName))
	}

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleting release %#q", releaseName))

		err := config.HelmClient.DeleteRelease(ctx, metav1.NamespaceDefault, releaseName, helmclient.DeleteOptions{})
		if err != nil {
			t.Fatalf("expected nil error got %v", err)
		}

		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deleted release %#q", releaseName))
	}
}
