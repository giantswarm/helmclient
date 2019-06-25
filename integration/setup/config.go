// +build k8srequired

package setup

import (
	"github.com/giantswarm/clientk8s"
	"github.com/giantswarm/helmclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	tillerNamespace = "giantswarm"
)

type Config struct {
	CPK8sClients *clientk8s.Clients
	HelmClient   *helmclient.Client
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

	var cpK8sClients *clientk8s.Clients
	{
		c := clientk8s.ClientsConfig{
			Logger: logger,

			KubeConfigPath: e2eHarnessDefaultKubeconfig,
		}

		cpK8sClients, err = clientk8s.NewClients(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var helmClient *helmclient.Client
	{
		c := helmclient.Config{
			K8sClient: cpK8sClients.K8sClient(),
			Logger:    logger,

			RestConfig:      cpK8sClients.RestConfig(),
			TillerNamespace: tillerNamespace,
		}

		helmClient, err = helmclient.New(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	c := Config{
		CPK8sClients: cpK8sClients,
		HelmClient:   helmClient,
		Logger:       logger,
	}

	return c, nil
}
