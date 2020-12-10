// +build k8srequired

package setup

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	{
		namespace := "giantswarm"
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}

		_, err = config.CPK8sClients.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
		if err != nil {
			return 1, microerror.Mask(err)
		}
	}

	return m.Run(), nil
}
