// +build k8srequired

package updatetiller

import (
	"context"
	"testing"

	"github.com/giantswarm/microerror"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUpdateTiller(t *testing.T) {
	ctx := context.Background()

	var err error

	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		t.Fatalf("could not install tiller %#v", err)
	}

	namespace := "giantswarm"
	labelSelector := "app=helm,name=tiller"
	outdatedTillerImage := "gcr.io/kubernetes-helm:v2.7.1"

	latestTillerImage, err := getTillerImage(namespace, labelSelector)
	if err != nil {
		t.Fatalf("could not get tiller image %#v", err)
	}

	err = updateTillerImage(namespace, labelSelector, outdatedTillerImage)
	if err != nil {
		t.Fatalf("could not set tiller image %#v", err)
	}

	downgradedTillerImage, err := getTillerImage(namespace, labelSelector)
	if err != nil {
		t.Fatalf("could not get tiller image %#v", err)
	}
	if downgradedTillerImage != outdatedTillerImage {
		t.Fatalf("tiller has not been downgraded got image %#q expected %#q", downgradedTillerImage, outdatedTillerImage)
	}

	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		t.Fatalf("could not install tiller %#v", err)
	}

	upgradedTillerImage, err := getTillerImage(namespace, labelSelector)
	if err != nil {
		t.Fatalf("could not get tiller image %#v", err)
	}

	if upgradedTillerImage != latestTillerImage {
		t.Fatalf("tiller has not been upgraded to latest image got %#q expected %#q", upgradedTillerImage, outdatedTillerImage)
	}
}

func getTillerDeployment(namespace string, labelSelector string) (*v1beta1.Deployment, error) {
	o := metav1.ListOptions{
		LabelSelector: labelSelector,
	}
	d, err := config.K8sClient.Extensions().Deployments(namespace).List(o)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if len(d.Items) > 1 {
		return nil, microerror.Maskf(tooManyResultsError, "%d", len(d.Items))
	}
	if len(d.Items) == 0 {
		return nil, microerror.Maskf(notFoundError, "%s", labelSelector)
	}

	return &d.Items[0], nil
}

func getTillerImage(namespace, labelSelector string) (string, error) {
	d, err := getTillerDeployment(namespace, labelSelector)
	if err != nil {
		return "", microerror.Mask(err)
	}

	if len(d.Spec.Template.Spec.Containers) > 1 {
		return "", microerror.Maskf(tooManyResultsError, "%d", len(d.Spec.Template.Spec.Containers))
	}
	if len(d.Spec.Template.Spec.Containers) == 0 {
		return "", microerror.Mask(notFoundError)
	}

	tillerImage := d.Spec.Template.Spec.Containers[0].Image
	if tillerImage == "" {
		return "", microerror.Maskf(notFoundError, "tiller image is empty")
	}

	return tillerImage, nil
}

func updateTillerImage(namespace, labelSelector, tillerImage string) error {
	deploy, err := getTillerDeployment(namespace, labelSelector)
	if err != nil {
		return microerror.Mask(err)
	}

	deploy.Spec.Template.Spec.Containers[0].Image = tillerImage
	_, err = config.K8sClient.Extensions().Deployments(namespace).Update(deploy)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
