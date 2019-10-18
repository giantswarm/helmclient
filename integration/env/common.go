// +build k8srequired

package env

import (
	"fmt"
	"os"
)

const (
	EnvVarCircleCI      = "CIRCLECI"
	EnvVarCircleSHA     = "CIRCLE_SHA1"
	EnvVarKeepResources = "KEEP_RESOURCES"
	EnvVarTestDir       = "TEST_DIR"

	// e2eHarnessDefaultKubeconfig is defined to avoid dependency of
	// e2e-harness. e2e-harness depends on this project. We don't want
	// circular dependencies even though it works in this case. This makes
	// vendoring very tricky.
	//
	// NOTE this should reflect value of DefaultKubeConfig constant.
	//
	//	See https://godoc.org/github.com/giantswarm/e2e-harness/pkg/harness#pkg-constants.
	//
	// There is also a note in the code there.
	//
	//	See https://github.com/giantswarm/e2e-harness/pull/177
	//
	e2eHarnessDefaultKubeconfig = "/workdir/.shipyard/config"
)

var (
	circleCI      string
	circleSHA     string
	keepResources string
	testDir       string
)

func init() {
	circleCI = os.Getenv(EnvVarCircleCI)
	circleSHA = os.Getenv(EnvVarCircleSHA)
	if circleSHA == "" {
		panic(fmt.Sprintf("env var %#q must not be empty", EnvVarCircleSHA))
	}

	keepResources = os.Getenv(EnvVarKeepResources)

	kubeconfig = os.Getenv(EnvVarE2EKubeconfig)
	if kubeconfig == "" {
		kubeconfig = e2eHarnessDefaultKubeconfig
	}
}

func CircleCI() bool {
	return circleCI == "true"
}

func CircleSHA() string {
	return circleSHA
}

func KeepResources() bool {
	return keepResources == "true"
}

func KubeConfigPath() string {
	return kubeconfig
}
