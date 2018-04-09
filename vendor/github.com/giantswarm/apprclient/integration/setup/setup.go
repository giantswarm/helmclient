// +build k8srequired

package setup

import (
	"testing"
)

func WrapTestMain(m *testing.M) {
	m.Run()
}
