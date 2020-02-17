// +build k8srequired

package basic

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/helmclient/integration/charttarball"
)

func TestBasic(t *testing.T) {
	ctx := context.Background()

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

	var releaseName string = "test-chart"

	{
		config.Logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("installing %#q", releaseName))

		tarballPath, err := charttarball.Create("test-chart")
		if err != nil {
			t.Fatalf("could not create chart archive %#v", err)
		}
		defer os.Remove(tarballPath)

		installOptions := helmclient.InstallOptions{
			ReleaseName: releaseName,
			Wait:        true,
		}
		err = config.HelmClient.InstallReleaseFromTarball(ctx, tarballPath, metav1.NamespaceDefault, map[string]interface{}{}, installOptions)
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
}
