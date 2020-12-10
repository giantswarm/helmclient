// +build k8srequired

package setup

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/giantswarm/helmclient/v3/integration/env"
	"github.com/giantswarm/helmclient/v3/pkg/helmclient"
)

type Config struct {
	CPK8sClients kubernetes.Interface
	HelmClient   helmclient.Interface
	Logger       micrologger.Logger
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

	config, err := clientcmd.BuildConfigFromFlags("", env.KubeConfigPath())
	if err != nil {
		return Config{}, microerror.Mask(err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return Config{}, microerror.Mask(err)
	}

	var helmClient *helmclient.Client
	{
		c := helmclient.Config{
			K8sClient:  k8sClient,
			Logger:     logger,
			RestClient: k8sClient.RESTClient(),
			RestConfig: config,
		}

		helmClient, err = helmclient.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	c := Config{
		CPK8sClients: k8sClient,
		HelmClient:   helmClient,
		Logger:       logger,
	}

	return c, nil
}
