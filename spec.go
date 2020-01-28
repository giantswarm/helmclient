package helmclient

import (
	"context"
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
	DeleteRelease(ctx context.Context, releaseName string) error
	// GetReleaseContent gets the current status of the Helm Release. The
	// releaseName is the name of the Helm Release that is set when the Chart
	// is installed.
	GetReleaseContent(ctx context.Context, releaseName string) (*ReleaseContent, error)
	// GetReleaseHistory gets the current installed version of the Helm Release.
	// The releaseName is the name of the Helm Release that is set when the Helm
	// Chart is installed.
	GetReleaseHistory(ctx context.Context, releaseName string) (*ReleaseHistory, error)
	// InstallReleaseFromTarball installs a Helm Chart packaged in the given tarball.
	InstallReleaseFromTarball(ctx context.Context, chartPath string, values map[string]interface{}, options InstallOptions) error
	// ListReleaseContents gets the current status of all Helm Releases.
	ListReleaseContents(ctx context.Context) ([]*ReleaseContent, error)
	// LoadChart loads a Helm Chart and returns its structure.
	LoadChart(ctx context.Context, chartPath string) (Chart, error)
	// PullChartTarball downloads a tarball from the provided tarball URL,
	// returning the file path.
	PullChartTarball(ctx context.Context, tarballURL string) (string, error)
	// RunReleaseTest runs the tests for a Helm Release. This is the same
	// action as running the helm test command.
	RunReleaseTest(ctx context.Context, releaseName string, options ReleaseTestOptions) error
	// UpdateReleaseFromTarball updates the given release using the chart packaged
	// in the tarball.
	UpdateReleaseFromTarball(ctx context.Context, releaseName, chartPath string, values map[string]interface{}, options UpdateOptions) error
}

// InstallOptions is the subset of supported options when installing helm
// releases.
type InstallOptions struct {
	Namespace   string
	ReleaseName string
	Wait        bool
}

// ReleaseTestOptions is the subset of supported options when running helm
// release tests.
type ReleaseTestOptions struct {
	Namespace string
}

// UpdateOptions is the subset of supported options when updating helm releases.
type UpdateOptions struct {
	Force     bool
	Namespace string
	Wait      bool
}
