package helmclient

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/k8sportforward"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/cmd/helm/installer"
	"k8s.io/helm/pkg/chartutil"
	helmclient "k8s.io/helm/pkg/helm"
)

// Config represents the configuration used to create a helm client.
type Config struct {
	// HelmClient sets a helm client used for all operations of the initiated
	// client. If this is nil, a new helm client will be created for each
	// operation via proper port forwarding. Setting the helm client here manually
	// might only be sufficient for testing or whenever you know what you do.
	HelmClient helmclient.Interface
	K8sClient  kubernetes.Interface
	Logger     micrologger.Logger

	RestConfig *rest.Config
}

// Client knows how to talk with a Helm Tiller server.
type Client struct {
	helmClient helmclient.Interface
	k8sClient  kubernetes.Interface
	logger     micrologger.Logger

	restConfig *rest.Config
}

// New creates a new configured Helm client.
func New(config Config) (*Client, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}

	if config.RestConfig == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.RestConfig must not be empty", config)
	}

	c := &Client{
		helmClient: config.HelmClient,
		k8sClient:  config.K8sClient,
		logger:     config.Logger,

		restConfig: config.RestConfig,
	}

	return c, nil
}

// DeleteRelease uninstalls a chart given its release name.
func (c *Client) DeleteRelease(releaseName string, options ...helmclient.DeleteOption) error {
	t, err := c.newTunnel()
	if err != nil {
		return microerror.Mask(err)
	}
	defer c.closeTunnel(t)

	_, err = c.newHelmClientFromTunnel(t).DeleteRelease(releaseName, options...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// GetReleaseContent gets the current status of the Helm Release including any
// values provided when the chart was installed. The releaseName is the name
// of the Helm Release that is set when the Helm Chart is installed.
func (c *Client) GetReleaseContent(releaseName string) (*ReleaseContent, error) {
	t, err := c.newTunnel()
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer c.closeTunnel(t)

	resp, err := c.newHelmClientFromTunnel(t).ReleaseContent(releaseName)
	if IsReleaseNotFound(err) {
		return nil, microerror.Maskf(releaseNotFoundError, releaseName)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	// If parameterizable values were passed at release creation time, raw values
	// are returned by the Tiller API and we convert these to a map. First we need
	// to check if there are values actually passed.
	var values chartutil.Values
	if resp.Release.Config != nil {
		raw := []byte(resp.Release.Config.Raw)
		values, err = chartutil.ReadValues(raw)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	content := &ReleaseContent{
		Name:   resp.Release.Name,
		Status: resp.Release.Info.Status.Code.String(),
		Values: values.AsMap(),
	}

	return content, nil
}

// GetReleaseHistory gets the current installed version of the Helm Release.
// The releaseName is the name of the Helm Release that is set when the Helm
// Chart is installed.
func (c *Client) GetReleaseHistory(releaseName string) (*ReleaseHistory, error) {
	t, err := c.newTunnel()
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer c.closeTunnel(t)

	resp, err := c.newHelmClientFromTunnel(t).ReleaseHistory(releaseName, helmclient.WithMaxHistory(1))
	if IsReleaseNotFound(err) {
		return nil, microerror.Maskf(releaseNotFoundError, releaseName)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}
	if len(resp.Releases) > 1 {
		return nil, microerror.Maskf(tooManyResultsError, "%d releases found, expected 1", len(resp.Releases))
	}

	var history *ReleaseHistory
	{
		release := resp.Releases[0]

		var version string
		if release.Chart != nil && release.Chart.Metadata != nil {
			version = release.Chart.Metadata.Version
		}

		history = &ReleaseHistory{
			Name:    release.Name,
			Version: version,
		}
	}

	return history, nil
}

// InstallFromTarball installs a chart packaged in the given tarball.
func (c *Client) InstallFromTarball(path, ns string, options ...helmclient.InstallOption) error {
	t, err := c.newTunnel()
	if err != nil {
		return microerror.Mask(err)
	}
	defer c.closeTunnel(t)

	_, err = c.newHelmClientFromTunnel(t).InstallRelease(path, ns, options...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) InstallTiller() error {
	var name = "tiller"
	var namespace = "kube-system"

	// Create the service account for tiller so it can pull images and do its do.
	{
		n := namespace
		i := &corev1.ServiceAccount{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ServiceAccount",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}

		_, err := c.k8sClient.CoreV1().ServiceAccounts(n).Create(i)
		if errors.IsAlreadyExists(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	// Create the cluster role binding for tiller so it is allowed to do its job.
	{
		i := &rbacv1.ClusterRoleBinding{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "rbac.authorization.k8s.io/v1",
				Kind:       "ClusterRoleBinding",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      name,
					Namespace: namespace,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "cluster-admin",
			},
		}

		_, err := c.k8sClient.RbacV1().ClusterRoleBindings().Create(i)
		if errors.IsAlreadyExists(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	// Install the tiller deployment in the guest cluster.
	{
		o := &installer.Options{
			ImageSpec:      "gcr.io/kubernetes-helm/tiller:v2.8.2",
			Namespace:      namespace,
			ServiceAccount: name,
		}

		err := installer.Install(c.k8sClient, o)
		if errors.IsAlreadyExists(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	// Wait for tiller to be up and running. When verifying to be able to ping
	// tiller we make sure 3 consecutive pings succeed before assuming everything
	// is fine.
	{
		c.logger.Log("level", "debug", "message", "attempt pinging tiller")

		var c int

		o := func() error {
			t, err := c.newTunnel()
			if err != nil {
				return microerror.Mask(err)
			}
			defer c.closeTunnel(t)

			err = c.newHelmClientFromTunnel(t).PingTiller()
			if err != nil {
				c = 0
				return microerror.Mask(err)
			}

			if c < 3 {
				return microerror.Maskf(executionFailedError, "failed pinging tiller")
			}
			c++

			return nil
		}
		b := newExponentialBackoff(2 * time.Minute)
		n := func(err error, delay time.Duration) {
			c.logger.Log("level", "debug", "message", "failed pinging tiller")
		}

		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return microerror.Mask(err)
		}

		c.logger.Log("level", "debug", "message", "succeeded pinging tiller")
	}

	return nil
}

// UpdateReleaseFromTarball updates the given release using the chart packaged
// in the tarball.
func (c *Client) UpdateReleaseFromTarball(releaseName, path string, options ...helmclient.UpdateOption) error {
	t, err := c.newTunnel()
	if err != nil {
		return microerror.Mask(err)
	}
	defer c.closeTunnel(t)

	_, err = c.newHelmClientFromTunnel(t).UpdateRelease(releaseName, path, options...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) closeTunnel(t *k8sportforward.Tunnel) {
	// In case a helm client is configured there is no tunnel and thus we do
	// nothing here.
	if t == nil {
		return
	}

	err := t.Close()
	if err != nil {
		c.logger.Log("level", "error", "message", "failed closing tunnel", "stack", fmt.Sprintf("%#v", err))
	}
}

func (c *Client) newHelmClientFromTunnel(t *k8sportforward.Tunnel) helmclient.Interface {
	// In case a helm client is configured we just go with it.
	if c.helmClient != nil {
		return c.helmClient
	}

	return helmclient.NewClient(
		helmclient.Host(newTunnelAddress(t)),
		helmclient.ConnectTimeout(5),
	)
}

func (c *Client) newTunnel() (*k8sportforward.Tunnel, error) {
	// In case a helm client is configured we do not need to create any port
	// forwarding.
	if c.helmClient != nil {
		return nil, nil
	}

	podName, err := getPodName(c.k8sClient, tillerLabelSelector, tillerDefaultNamespace)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var forwarder *k8sportforward.Forwarder
	{
		c := k8sportforward.Config{
			RestConfig: c.restConfig,
		}
		forwarder, err = k8sportforward.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var tunnel *k8sportforward.Tunnel
	{
		c := k8sportforward.TunnelConfig{
			Remote:    tillerPort,
			Namespace: tillerDefaultNamespace,
			PodName:   podName,
		}

		tunnel, err = forwarder.ForwardPort(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return tunnel, nil
}

func getPodName(client kubernetes.Interface, labelSelector, namespace string) (string, error) {
	o := metav1.ListOptions{
		LabelSelector: labelSelector,
	}
	pods, err := client.CoreV1().Pods(namespace).List(o)
	if err != nil {
		return "", microerror.Mask(err)
	}

	if len(pods.Items) > 1 {
		return "", microerror.Mask(tooManyResultsError)
	}
	if len(pods.Items) == 0 {
		return "", microerror.Mask(notFoundError)
	}
	pod := pods.Items[0]

	return pod.Name, nil
}

// TODO remove when k8sportforward.Tunnel.Address() got implemented.
func newTunnelAddress(t *k8sportforward.Tunnel) string {
	return fmt.Sprintf("127.0.0.1:%d", t.Local)
}
