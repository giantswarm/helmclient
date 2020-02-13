// +build k8srequired

package basic

import (
	"context"
	"testing"
)

func TestBasic(t *testing.T) {
	ctx := context.Background()
	var err error

	releaseContent, err := config.HelmClient.GetReleaseContent(ctx, "no-release-exists")
	if err != nil {
		t.Fatalf("could not get release content %v", err)
	}
	if releaseContent != nil {
		t.Fatalf("expected nil release got %v", err)
	}
}
