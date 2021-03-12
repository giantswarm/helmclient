module github.com/giantswarm/helmclient/v4

go 1.14

require (
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/kubeconfig/v4 v4.0.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/micrologger v0.5.0
	github.com/go-logr/logr v0.3.0 // indirect
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/google/go-cmp v0.5.5
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/mholt/archiver/v3 v3.5.0
	github.com/moby/term v0.0.0-20201110203204-bea5bbe245bf // indirect
	github.com/onsi/ginkgo v1.14.1 // indirect
	github.com/onsi/gomega v1.10.2 // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/afero v1.5.1
	go.uber.org/zap v1.15.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	helm.sh/helm/v3 v3.5.3
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.6.5
)

replace (
	github.com/containerd/containerd v1.3.4 => github.com/containerd/containerd v1.4.4
	github.com/coreos/etcd v3.3.10+incompatible => github.com/coreos/etcd v3.3.25+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	// Use moby v20.10.5 to fix build issue on darwin.
	github.com/docker/docker => github.com/moby/moby v20.10.5+incompatible
	github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.4.2
	// Use mergo 0.3.11 due to bug in 0.3.9 merging Go structs.
	github.com/imdario/mergo => github.com/imdario/mergo v0.3.11
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc93
	// Use fork of CAPI with Kubernetes 1.18 support.
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.13-gs
)
