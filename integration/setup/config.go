// +build k8srequired

package setup

import (
	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/giantswarm/helmclient"
)

const (
	tillerNamespace = "giantswarm"
)

type Config struct {
	HelmClient helmclient.Interface
	Host       *framework.Host
	K8sClient  kubernetes.Interface
	Logger     micrologger.Logger

	RestConfig *rest.Config
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

	var helmClient helmclient.Interface
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

	var host *framework.Host
	{
		c := framework.HostConfig{
			Logger: logger,

			ClusterID:  "na",
			VaultToken: "na",
		}

		host, err = framework.NewHost(c)
		if err != nil {
			return Config{}, microerror.Mask(err)
		}
	}

	c := Config{
		HelmClient: helmClient,
		Host:       host,
		K8sClient:  k8sClient,
		Logger:     logger,

		RestConfig: restConfig,
	}

	return c, nil
}
