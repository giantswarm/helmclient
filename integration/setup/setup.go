// +build k8srequired

package setup

import (
	"log"
	"os"
	"testing"

	"github.com/giantswarm/e2e-harness/pkg/framework"
)

func WrapTestMain(h *framework.Host, m *testing.M) {
	var v int
	var err error

	if err = h.CreateNamespace("giantswarm"); err != nil {
		log.Printf("%#v\n", err)
		v = 1
	}

	if v == 0 {
		m.Run()
	}

	if os.Getenv("KEEP_RESOURCES") != "true" {
		// only do full teardown when not on CI
		if os.Getenv("CIRCLECI") != "true" {
			// TODO there should be error handling for the framework teardown.
			h.Teardown()
		}
	}

	os.Exit(v)
}
