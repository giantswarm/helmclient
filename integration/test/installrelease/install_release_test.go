// +build k8srequired

package basic

import (
	"context"
	"os"
	"testing"

	"k8s.io/helm/pkg/helm"

	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/helmclient/integration/charttarball"
)

func TestInstallRelease_IsReleaseAlreadyExists(t *testing.T) {
	ctx := context.Background()
	var err error

	const releaseName = "test"

	tarballPath, err := charttarball.Create("test-chart")
	if err != nil {
		t.Fatalf("could not create chart archive %#v", err)
	}
	defer os.Remove(tarballPath)

	// We need to pass the ValueOverrides option to make the install process
	// use the default values and prevent errors on nested values.
	//
	//	{
	//		rpc error: code = Unknown desc = render error in "cnr-server-chart/templates/deployment.yaml":
	//		template: cnr-server-chart/templates/deployment.yaml:20:26:
	//		executing "cnr-server-chart/templates/deployment.yaml" at <.Values.image.reposi...>: can't evaluate field repository in type interface {}
	//	}
	//
	err = config.HelmClient.InstallReleaseFromTarball(ctx, tarballPath, map[string]interface{}{}, InstallOptions{})
	if err != nil {
		t.Fatalf("failed to install release %#v", err)
	}

	err = config.HelmClient.InstallReleaseFromTarball(ctx, tarballPath, values, InstallOptions{})
	if helmclient.IsReleaseAlreadyExists(err) {
		// This is error we want.
	} else if err != nil {
		t.Fatalf("failed to install release %#v", err)
	} else {
		t.Fatalf("expected error for already installed release, got %#v", err)
	}
}

func TestInstallRelease_IsTarballNotFound(t *testing.T) {
	ctx := context.Background()
	var err error

	const releaseName = "test"
	const tarballPath = "/path/that/does/not-exist"

	// We need to pass the ValueOverrides option to make the install process
	// use the default values and prevent errors on nested values.
	//
	//	{
	//		rpc error: code = Unknown desc = render error in "cnr-server-chart/templates/deployment.yaml":
	//		template: cnr-server-chart/templates/deployment.yaml:20:26:
	//		executing "cnr-server-chart/templates/deployment.yaml" at <.Values.image.reposi...>: can't evaluate field repository in type interface {}
	//	}
	//
	err = config.HelmClient.InstallReleaseFromTarball(ctx, tarballPath, map[string]interface{}{}, InstallOptions{})
	if helmclient.IsTarballNotFound(err) {
		// This is error we want.
	} else if err != nil {
		t.Fatalf("failed to install release %#v", err)
	} else {
		t.Fatalf("expected error on missing tarball, got %#v", err)
	}
}
