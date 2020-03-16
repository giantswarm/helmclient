package helmclient

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Describes the status of a release. This needs to be kept in sync with
// upstream but it allows us to have constants without importing Helm
// packages.
//
// See: https://github.com/helm/helm/blob/master/pkg/release/status.go
const (
	// StatusUnknown indicates that a release is in an uncertain state.
	StatusUnknown = "unknown"
	// StatusDeployed indicates that the release has been pushed to Kubernetes.
	StatusDeployed = "deployed"
	// StatusUninstalled indicates that a release has been uninstalled from Kubernetes.
	StatusUninstalled = "uninstalled"
	// StatusSuperseded indicates that this release object is outdated and a newer one exists.
	StatusSuperseded = "superseded"
	// StatusFailed indicates that the release was not successfully deployed.
	StatusFailed = "failed"
	// StatusUninstalling indicates that a uninstall operation is underway.
	StatusUninstalling = "uninstalling"
	// StatusPendingInstall indicates that an install operation is underway.
	StatusPendingInstall = "pending-install"
	// StatusPendingUpgrade indicates that an upgrade operation is underway.
	StatusPendingUpgrade = "pending-upgrade"
	// StatusPendingRollback indicates that an rollback operation is underway.
	StatusPendingRollback = "pending-rollback"
)

var (
	// ReleaseTransitionStatuses is used to determine if the Helm Release is
	// currently being updated.
	ReleaseTransitionStatuses = map[string]bool{
		StatusUninstalled:     true,
		StatusPendingInstall:  true,
		StatusPendingUpgrade:  true,
		StatusPendingRollback: true,
	}
)

const (
	// defaultMaxHistory is the maximum number of release versions stored per
	// release by default.
	defaultMaxHistory = 10
	// httpClientTimeout is the timeout when pulling tarballs.
	httpClientTimeout = 5
	// runReleaseTestTimeout is the timeout in seconds when running tests.
	runReleaseTestTimout = 300
)

// Interface describes the methods provided by the helm client.
type Interface interface {
	// DeleteRelease uninstalls a chart given its release name.
	DeleteRelease(ctx context.Context, namespace, releaseName string) error
	// GetReleaseContent gets the current status of the Helm Release. The
	// releaseName is the name of the Helm Release that is set when the Chart
	// is installed.
	GetReleaseContent(ctx context.Context, namespace, releaseName string) (*ReleaseContent, error)
	// GetReleaseHistory gets the current installed version of the Helm Release.
	// The releaseName is the name of the Helm Release that is set when the Helm
	// Chart is installed.
	GetReleaseHistory(ctx context.Context, namespace, releaseName string) (*ReleaseHistory, error)
	// InstallReleaseFromTarball installs a Helm Chart packaged in the given tarball.
	InstallReleaseFromTarball(ctx context.Context, chartPath, namespace string, values map[string]interface{}, options InstallOptions) error
	// ListReleaseContents gets the current status of all Helm Releases.
	ListReleaseContents(ctx context.Context, namespace string) ([]*ReleaseContent, error)
	// LoadChart loads a Helm Chart and returns its structure.
	LoadChart(ctx context.Context, chartPath string) (Chart, error)
	// PullChartTarball downloads a tarball from the provided tarball URL,
	// returning the file path.
	PullChartTarball(ctx context.Context, tarballURL string) (string, error)
	// RunReleaseTest runs the tests for a Helm Release. This is the same
	// action as running the helm test command.
	RunReleaseTest(ctx context.Context, namespace, releaseName string) error
	// UpdateReleaseFromTarball updates the given release using the chart packaged
	// in the tarball.
	UpdateReleaseFromTarball(ctx context.Context, chartPath, namespace, releaseName string, values map[string]interface{}, options UpdateOptions) error
}

type RESTClientGetter interface {
	// ToDiscoveryClient returns discovery client
	ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error)
	// ToRawKubeConfigLoader return kubeconfig loader as-is
	ToRawKubeConfigLoader() clientcmd.ClientConfig
	// ToRESTConfig returns restconfig
	ToRESTConfig() (*rest.Config, error)
	// ToRESTMapper returns a restmapper
	ToRESTMapper() (meta.RESTMapper, error)
}

// InstallOptions is the subset of supported options when installing helm
// releases.
type InstallOptions struct {
	Namespace   string
	ReleaseName string
	Wait        bool
}

// UpdateOptions is the subset of supported options when updating helm releases.
type UpdateOptions struct {
	Force bool
	Wait  bool
}
