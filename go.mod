module github.com/giantswarm/helmclient/v4

go 1.14

require (
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/kubeconfig/v4 v4.0.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/micrologger v0.5.0
	github.com/go-logr/logr v0.3.0 // indirect
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/google/go-cmp v0.5.4
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/mholt/archiver v2.1.0+incompatible
	github.com/moby/term v0.0.0-20201110203204-bea5bbe245bf // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/onsi/ginkgo v1.14.1 // indirect
	github.com/onsi/gomega v1.10.2 // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/afero v1.5.1
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/appengine v1.6.6 // indirect
	gotest.tools/v3 v3.0.3 // indirect
	helm.sh/helm/v3 v3.4.2
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v0.19.4
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/controller-runtime v0.6.5
)

replace (
	// Use v0.8.10 of hcsshim to fix nancy alert.
	github.com/Microsoft/hcsshim v0.8.7 => github.com/Microsoft/hcsshim v0.8.10
	// Use moby v20.10.0-beta1 to fix build issue on darwin.
	github.com/docker/docker => github.com/moby/moby v20.10.0-beta1+incompatible
	// Use mergo 0.3.11 due to bug in 0.3.9 merging Go structs.
	github.com/imdario/mergo => github.com/imdario/mergo v0.3.11
	// Use v1.0.0-rc7 of runc to fix nancy alert.
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc7
	// Use fork of CAPI with Kubernetes 1.18 support.
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
)
