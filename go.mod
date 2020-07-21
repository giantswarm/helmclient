module github.com/giantswarm/helmclient

go 1.14

require (
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/k8sclient/v3 v3.1.2
	github.com/giantswarm/kubeconfig v0.2.1
	github.com/giantswarm/microerror v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/google/go-cmp v0.5.0
	github.com/mholt/archiver v2.1.0+incompatible
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/afero v1.3.2
	github.com/ulikunitz/xz v0.5.7 // indirect
	gopkg.in/yaml.v2 v2.3.0
	helm.sh/helm/v3 v3.1.3
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/controller-runtime v0.4.0
)
