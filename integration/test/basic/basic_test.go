// +build k8srequired

package basic

import (
	"path/filepath"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/helm/pkg/helm"
)

func TestInstallChart(t *testing.T) {
	l, err := micrologger.New(micrologger.Config{})
	if err != nil {
		t.Fatalf("could not create logger %v", err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", harness.DefaultKubeConfig)
	if err != nil {
		t.Fatalf("could not create k8s config %v", err)
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("could not create k8s client %v", err)
	}

	c := helmclient.Config{
		Logger:          l,
		K8sClient:       cs,
		RestConfig:      config,
		TillerNamespace: "giantswarm",
	}

	helmClient, err := helmclient.New(c)
	if err != nil {
		t.Fatalf("could not create helm client %v", err)
	}

	err = helmClient.EnsureTillerInstalled()
	if err != nil {
		t.Fatalf("could not install Tiller %v", err)
	}

	// --test-dir dir is mounted in /e2e in the test container.
	tarballPath := filepath.Join("/e2e/fixtures/", "tb-chart.tar.gz")

	const releaseName = "tb-chart-release"

	// We need to pass the ValueOverrides option to make the install process
	// use the default values and prevent errors on nested values.
	//
	//     {
	//      rpc error: code = Unknown desc = render error in "cnr-server-chart/templates/deployment.yaml":
	//      template: cnr-server-chart/templates/deployment.yaml:20:26:
	//      executing "cnr-server-chart/templates/deployment.yaml" at <.Values.image.reposi...>: can't evaluate field repository in type interface {}
	//     }
	//
	err = helmClient.InstallFromTarball(tarballPath, "default", helm.ReleaseName(releaseName), helm.ValueOverrides([]byte("{}")))
	if err != nil {
		t.Fatalf("could not install chart %v", err)
	}

	releaseContent, err := helmClient.GetReleaseContent(releaseName)
	if err != nil {
		t.Fatalf("could not get release content %v", err)
	}

	expectedName := releaseName
	actualName := releaseContent.Name
	if expectedName != actualName {
		t.Fatalf("bad release name, want %q, got %q", expectedName, actualName)
	}

	expectedStatus := "DEPLOYED"
	actualStatus := releaseContent.Status
	if expectedStatus != actualStatus {
		t.Fatalf("bad release status, want %q, got %q", expectedStatus, actualStatus)
	}

	err = helmClient.RunReleaseTest(releaseName)
	if err != nil {
		t.Fatalf("error running tests, want nil got %v", err)
	}

	// Test should fail on the 2nd attempt because the test pod already exists.
	err = helmClient.RunReleaseTest(releaseName)
	if err == nil {
		t.Fatalf("error running tests, want error got nil")
	}

	err = helmClient.DeleteRelease(releaseName)
	if err != nil {
		t.Fatalf("could not delete release %v", err)
	}

	releaseContent, err = helmClient.GetReleaseContent(releaseName)
	if err != nil {
		t.Fatalf("could not get release content %v", err)
	}
	expectedStatus = "DELETED"
	actualStatus = releaseContent.Status
	if expectedStatus != actualStatus {
		t.Fatalf("bad release status, want %q, got %q", expectedStatus, actualStatus)
	}
}
