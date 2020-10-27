module github.com/giantswarm/helmclient/v3

go 1.14

require (
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/k8sclient/v5 v5.0.0
	github.com/giantswarm/kubeconfig/v3 v3.0.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.3
	github.com/google/go-cmp v0.5.2
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/mholt/archiver v2.1.0+incompatible
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/prometheus/client_golang v1.8.0
	github.com/spf13/afero v1.4.1
	github.com/ulikunitz/xz v0.5.7 // indirect
	gopkg.in/yaml.v2 v2.3.0
	helm.sh/helm/v3 v3.4.0
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/controller-runtime v0.6.3
)

replace sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
