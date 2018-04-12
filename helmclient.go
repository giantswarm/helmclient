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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/cmd/helm/installer"
	"k8s.io/helm/pkg/chartutil"
	helmclient "k8s.io/helm/pkg/helm"
)

const (
	connectionTimeoutSecs = 5
)

// Config represents the configuration used to create a helm client.
type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	RestConfig *rest.Config
}

// Client knows how to talk with a Helm Tiller server.
type Client struct {
	helmClient helmclient.Interface
	k8sClient  kubernetes.Interface
	logger     micrologger.Logger
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

	host, err := setupConnection(config.K8sClient, config.RestConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	helmClient := helmclient.NewClient(helmclient.Host(host), helmclient.ConnectTimeout(connectionTimeoutSecs))

	c := &Client{
		helmClient: helmClient,
		k8sClient:  config.K8sClient,
		logger:     config.Logger,
	}

	return c, nil
}

// DeleteRelease uninstalls a chart given its release name.
func (c *Client) DeleteRelease(releaseName string, options ...helmclient.DeleteOption) error {
	_, err := c.helmClient.DeleteRelease(releaseName, options...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// GetReleaseContent gets the current status of the Helm Release including any
// values provided when the chart was installed. The releaseName is the name
// of the Helm Release that is set when the Helm Chart is installed.
func (c *Client) GetReleaseContent(releaseName string) (*ReleaseContent, error) {
	resp, err := c.helmClient.ReleaseContent(releaseName)
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
	var version string

	resp, err := c.helmClient.ReleaseHistory(releaseName, helmclient.WithMaxHistory(1))
	if IsReleaseNotFound(err) {
		return nil, microerror.Maskf(releaseNotFoundError, releaseName)
	} else if err != nil {
		return nil, microerror.Mask(err)
	}
	if len(resp.Releases) > 1 {
		return nil, microerror.Maskf(tooManyResultsError, "%d releases found, expected 1", len(resp.Releases))
	}

	release := resp.Releases[0]
	if release.Chart != nil && release.Chart.Metadata != nil {
		version = release.Chart.Metadata.Version
	}

	history := &ReleaseHistory{
		Name:    release.Name,
		Version: version,
	}

	return history, nil
}

// InstallFromTarball installs a chart packaged in the given tarball.
func (c *Client) InstallFromTarball(path, ns string, options ...helmclient.InstallOption) error {
	_, err := c.helmClient.InstallRelease(path, ns, options...)
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
		if err != nil {
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
					APIGroup:  "rbac.authorization.k8s.io",
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
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// Install the tiller deployment in the guest cluster.
	{
		o := &installer.Options{
			Namespace:      namespace,
			ServiceAccount: name,
		}

		err := installer.Install(c.k8sClient, o)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// Wait for tiller to be up and running.
	{
		c.logger.Log("level", "debug", "message", "attempt pinging tiller")

		o := func() error {
			err := c.helmClient.PingTiller()
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}
		b := newExponentialBackoff(60 * time.Second)
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
	_, err := c.helmClient.UpdateRelease(releaseName, path, options...)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func setupConnection(client kubernetes.Interface, config *rest.Config) (string, error) {
	podName, err := getPodName(client, tillerLabelSelector, tillerDefaultNamespace)
	if err != nil {
		return "", microerror.Mask(err)
	}

	c := k8sportforward.Config{
		RestConfig: config,
	}
	f, err := k8sportforward.New(c)
	if err != nil {
		return "", microerror.Mask(err)
	}

	tc := k8sportforward.TunnelConfig{
		Remote:    tillerPort,
		Namespace: tillerDefaultNamespace,
		PodName:   podName,
	}

	t, err := f.ForwardPort(tc)
	if err != nil {
		return "", microerror.Mask(err)
	}

	host := fmt.Sprintf("127.0.0.1:%d", t.Local)

	return host, nil
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
