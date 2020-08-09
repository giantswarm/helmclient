// +build k8srequired

package setup

import (
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/helmclient/v2/integration/env"
	"github.com/giantswarm/helmclient/v2/pkg/helmclient"
)

type Config struct {
	CPK8sClients *k8sclient.Clients
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

	var cpK8sClients *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			Logger: logger,

			KubeConfigPath: env.KubeConfigPath(),
		}

		cpK8sClients, err = k8sclient.NewClients(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	var helmClient *helmclient.Client
	{
		c := helmclient.Config{
			K8sClient: cpK8sClients,
			Logger:    logger,
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
