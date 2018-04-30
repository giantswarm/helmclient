// +build k8srequired

package basic

import (
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
	"github.com/giantswarm/helmclient/integration/setup"
)

var (
	h *framework.Host
)

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	var err error

	h, err = framework.NewHost(framework.HostConfig{})
	if err != nil {
		panic(err.Error())
	}

	setup.WrapTestMain(h, m)
}
