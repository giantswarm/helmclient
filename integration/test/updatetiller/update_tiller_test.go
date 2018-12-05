// +build k8srequired

package updatetiller

import (
	"context"
	"testing"
)

func TestUpdateTiller(t *testing.T) {
	ctx := context.Background()
	var err error

	err = config.HelmClient.EnsureTillerInstalled(ctx)
	if err != nil {
		t.Fatalf("could not install Tiller %#v", err)
	}
}
