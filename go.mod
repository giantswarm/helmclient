module github.com/giantswarm/helmclient/v4

go 1.16

require (
	github.com/giantswarm/backoff v1.0.0
	github.com/giantswarm/kubeconfig/v4 v4.1.0
	github.com/giantswarm/microerror v0.4.0
	github.com/giantswarm/micrologger v0.6.0
	github.com/google/go-cmp v0.5.6
	github.com/mholt/archiver/v3 v3.5.0
	github.com/opencontainers/image-spec v1.0.2
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/afero v1.6.0
	helm.sh/helm/v3 v3.8.1
	k8s.io/api v0.23.4
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
	oras.land/oras-go v1.1.0
	sigs.k8s.io/controller-runtime v0.9.7
)

replace (
	github.com/containerd/containerd v1.3.4 => github.com/containerd/containerd v1.4.4
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v20.10.11+incompatible
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.4.2
	// Use mergo 0.3.11 due to bug in 0.3.9 merging Go structs.
	github.com/imdario/mergo => github.com/imdario/mergo v0.3.11
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc93
	github.com/ulikunitz/xz => github.com/ulikunitz/xz v0.5.10
)
