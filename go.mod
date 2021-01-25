module github.com/giantswarm/helmclient/v4

go 1.14

require (
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/kubeconfig/v4 v4.0.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/micrologger v0.5.0
	github.com/google/go-cmp v0.5.4
	github.com/mholt/archiver v2.1.0+incompatible
	github.com/moby/term v0.0.0-20201110203204-bea5bbe245bf // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/afero v1.5.1
	gotest.tools/v3 v3.0.3 // indirect
	helm.sh/helm/v3 v3.4.2
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/controller-runtime v0.8.1
)

replace (
	// Use moby v20.10.0-beta1 to fix build issue on darwin.
	github.com/docker/docker => github.com/moby/moby v20.10.0-beta1+incompatible
	// Use mergo 0.3.11 due to bug in 0.3.9 merging Go structs.
	github.com/imdario/mergo => github.com/imdario/mergo v0.3.11
	// Use fork of CAPI with Kubernetes 1.18 support.
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
)
