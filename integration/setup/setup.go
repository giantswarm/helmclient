// +build k8srequired

package setup

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/helmclient/integration/env"
	"github.com/giantswarm/microerror"
)

func Setup(m *testing.M, config Config) {
	ctx := context.Background()

	exitCode, err := setup(ctx, m, config)
	if err != nil {
		config.Logger.LogCtx(ctx, "level", "error", "message", "", "stack", fmt.Sprintf("%#v", err))
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func setup(ctx context.Context, m *testing.M, config Config) (int, error) {
	var err error
	teardown := !env.CircleCI() && !env.KeepResources()

	{
		// TODO this should be named EnsureNamespaceCreated
		err = config.K8sSetup.EnsureNamespace(ctx, tillerNamespace)
		if err != nil {
			return 1, microerror.Mask(err)
		}
		if teardown {
			// TODO defer config.K8sSetup.EnsureNamespaceDeleted(ctx, tillerNamespace)
		}
	}

	return m.Run(), nil
}
