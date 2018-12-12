// +build k8srequired

package upgradetiller

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUpgradeTiller(t *testing.T) {
	ctx := context.Background()

	var err error

	labelSelector := "app=helm,name=tiller"
	outdatedTillerImage := "gcr.io/kubernetes-helm:v2.7.2"
	tillerNamespace := "giantswarm"

	var helmClient *helmclient.Client
	{
		c := helmclient.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,

			RestConfig: config.RestConfig,
			// Use outdated tiller for initial install.
			TillerImage:     outdatedTillerImage,
			TillerNamespace: tillerNamespace,
		}

		helmClient, err = helmclient.New(c)
		if err != nil {
			t.Fatalf("could not create tiller client %#v", err)
		}
	}

	// Install tiller using helm client with outdated image.
	err = helmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		t.Fatalf("could not install tiller %#v", err)
	}

	tillerImage, err := getTillerImage(ctx, tillerNamespace, labelSelector)
	if err != nil {
		t.Fatalf("could not get tiller image %#v", err)
	}
	if tillerImage != outdatedTillerImage {
		t.Fatalf("tiller has not been downgraded got image %#q expected %#q", tillerImage, outdatedTillerImage)
	}

	// Upgrade tiller to the latest image using default helm client.
	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		t.Fatalf("could not install tiller %#v", err)
	}

	upgradedTillerImage, err := getTillerImage(ctx, tillerNamespace, labelSelector)
	if err != nil {
		t.Fatalf("could not get tiller image %#v", err)
	}

	if upgradedTillerImage != helmclient.TillerImageSpec {
		t.Fatalf("tiller has not been upgraded to latest image got %#q expected %#q", upgradedTillerImage, helmclient.TillerImageSpec)
	}
}

func getTillerDeployment(ctx context.Context, namespace string, labelSelector string) (*appsv1.Deployment, error) {
	var d *appsv1.Deployment
	{
		o := func() error {
			lo := metav1.ListOptions{
				LabelSelector: labelSelector,
			}
			l, err := config.K8sClient.Apps().Deployments(namespace).List(lo)
			if err != nil {
				return microerror.Mask(err)
			}

			if len(l.Items) != 1 {
				return microerror.Maskf(executionFailedError, "cannot get deployment for %#q %#q found %d, want 1", namespace, labelSelector, len(l.Items))
			}

			d = &l.Items[0]
			if d.Status.AvailableReplicas != 1 && d.Status.ReadyReplicas != 1 {
				return microerror.Maskf(executionFailedError, "tiller deployment not ready %d available %d ready, want 1", d.Status.AvailableReplicas, d.Status.ReadyReplicas)
			}

			return nil
		}

		b := backoff.NewExponential(2*time.Minute, 5*time.Second)
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

	if len(d.Spec.Template.Spec.Containers) != 1 {
		return "", microerror.Maskf(executionFailedError, "Spec.Template.Spec.Containers == %d, want 1", len(d.Spec.Template.Spec.Containers))
	}

	tillerImage := d.Spec.Template.Spec.Containers[0].Image
	if tillerImage == "" {
		return "", microerror.Maskf(executionFailedError, "tiller image is empty")
	}

	return tillerImage, nil
}
