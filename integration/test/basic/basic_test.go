// +build k8srequired

package basic

import (
	"context"
	"testing"

	"github.com/giantswarm/helmclient"
)

func TestBasic(t *testing.T) {
	ctx := context.Background()
	var err error

	releaseContent, err := config.HelmClient.GetReleaseContent(ctx, "no-release-exists", "kube-system")
	if err != nil && !helmclient.IsReleaseNotFound(err) {
		t.Fatalf("expected release not found error got %v", err)
	}
	if releaseContent != nil {
		t.Fatalf("expected nil release got %v", err)
	}
}
