// +build k8srequired

package basic

import (
	"testing"

	"github.com/giantswarm/apprclient/integration/setup"
)

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	setup.WrapTestMain(m)
}
