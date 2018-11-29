// +build k8srequired

package setup

import (
	"github.com/giantswarm/e2esetup/k8s"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	tillerNamespace = "giantswarm"
)

type Config struct {
	HelmClient *helmclient.Client
	K8sSetup   *k8s.Setup
	Logger     micrologger.Logger
}

func NewConfig() (Config, error) {
	var err error

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var restConfig *rest.Config
	{
		restConfig, err = clientcmd.BuildConfigFromFlags("", e2eHarnessDefaultKubeconfig)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var k8sClient *kubernetes.Clientset
	{
		k8sClient, err = kubernetes.NewForConfig(restConfig)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var k8sSetup *k8s.Setup
	{
		c := k8s.SetupConfig{
			K8sClient: k8sClient,
			Logger:    logger,
		}

		k8sSetup, err = k8s.NewSetup(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var helmClient *helmclient.Client
	{
		c := helmclient.Config{
			K8sClient: k8sClient,
			Logger:    logger,

			RestConfig:      restConfig,
			TillerNamespace: tillerNamespace,
		}

		helmClient, err = helmclient.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	c := Config{
		HelmClient: helmClient,
		K8sSetup:   k8sSetup,
		Logger:     logger,
	}

	return c, nil
}
