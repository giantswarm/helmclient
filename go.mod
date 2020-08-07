module github.com/giantswarm/helmclient/v2

go 1.14

require (
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/helmclient v1.0.6 // indirect
	github.com/giantswarm/k8sclient/v4 v4.0.0-20200806115259-2d3b230ace59
	github.com/giantswarm/kubeconfig/v2 v2.0.0-20200806114529-9ef53912cf03
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.1
	github.com/google/go-cmp v0.5.1
	github.com/mholt/archiver v2.1.0+incompatible
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/afero v1.3.3
	gopkg.in/yaml.v2 v2.3.0
	helm.sh/helm/v3 v3.2.4
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.5
	sigs.k8s.io/controller-runtime v0.6.1
)
