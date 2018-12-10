// +build k8srequired

package upgradetiller

import (
	"context"
	"testing"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUpgradeTiller(t *testing.T) {
	ctx := context.Background()

	var err error

	// Install tiller with the latest image.
	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		t.Fatalf("could not install tiller %#v", err)
	}

	namespace := "giantswarm"
	labelSelector := "app=helm,name=tiller"
	outdatedTillerImage := "gcr.io/kubernetes-helm:v2.7.1"

	latestTillerImage, err := getTillerImage(ctx, namespace, labelSelector)
	if err != nil {
		t.Fatalf("could not get tiller image %#v", err)
	}

	// Downgrade tiller to a previous version.
	err = updateTillerImage(ctx, namespace, labelSelector, outdatedTillerImage)
	if err != nil {
		t.Fatalf("could not set tiller image %#v", err)
	}

	downgradedTillerImage, err := getTillerImage(ctx, namespace, labelSelector)
	if err != nil {
		t.Fatalf("could not get tiller image %#v", err)
	}
	if downgradedTillerImage != outdatedTillerImage {
		t.Fatalf("tiller has not been downgraded got image %#q expected %#q", downgradedTillerImage, outdatedTillerImage)
	}

	// Upgrade tiller to the latest image.
	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		t.Fatalf("could not install tiller %#v", err)
	}

	upgradedTillerImage, err := getTillerImage(ctx, namespace, labelSelector)
	if err != nil {
		t.Fatalf("could not get tiller image %#v", err)
	}

	if upgradedTillerImage != latestTillerImage {
		t.Fatalf("tiller has not been upgraded to latest image got %#q expected %#q", upgradedTillerImage, outdatedTillerImage)
	}
}

func getTillerDeployment(ctx context.Context, namespace string, labelSelector string) (*appsv1.Deployment, error) {
	var d *appsv1.Deployment
	{
		o := func() error {
			lo := metav1.ListOptions{
				LabelSelector: labelSelector,
			}
			dl, err := config.K8sClient.Apps().Deployments(namespace).List(lo)
			if err != nil {
				return microerror.Mask(err)
			}

			if len(dl.Items) > 1 {
				return microerror.Maskf(tooManyResultsError, "%d", len(dl.Items))
			}
			if len(dl.Items) == 0 {
				return microerror.Maskf(notFoundError, "%s", labelSelector)
			}

			d = &dl.Items[0]
			if d.Status.AvailableReplicas != 1 && d.Status.ReadyReplicas != 1 {
				return microerror.Maskf(notFoundError, "tiller deployment updating expected 1 pod found %d available %d ready", d.Status.AvailableReplicas, d.Status.ReadyReplicas)
			}

			return nil
		}

		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		n := backoff.NewNotifier(config.Logger, ctx)

		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return d, nil
}

func getTillerImage(ctx context.Context, namespace, labelSelector string) (string, error) {
	d, err := getTillerDeployment(ctx, namespace, labelSelector)
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

func updateTillerImage(ctx context.Context, namespace, labelSelector, tillerImage string) error {
	d, err := getTillerDeployment(ctx, namespace, labelSelector)
	if err != nil {
		return microerror.Mask(err)
	}

	d.Spec.Template.Spec.Containers[0].Image = tillerImage
	_, err = config.K8sClient.Apps().Deployments(namespace).Update(d)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
