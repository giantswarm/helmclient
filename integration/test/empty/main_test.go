// +build k8srequired

package empty

import (
	"fmt"
	"testing"
)

// TestMain allows us to have common setup and teardown steps that are run
// once for all the tests https://golang.org/pkg/testing/#hdr-Main.
func TestMain(m *testing.M) {
	fmt.Println("Empty test for debugging orb")
}
