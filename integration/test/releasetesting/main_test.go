//go:build k8srequired
// +build k8srequired

package releasetesting

import (
	"fmt"
	"testing"

	"github.com/giantswarm/helmclient/v4/integration/setup"
)

var (
	config setup.Config
)

func init() {
	var err error

	{
		config, err = setup.NewConfig()
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}
}

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	setup.Setup(m, config)
}
