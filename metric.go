package helmclient

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	operationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "helmclient",
			Name:      "operation",
			Help:      "Histogram of operations carried out by the helm client.",
		},
		[]string{"operation"},
	)
)

func init() {
	prometheus.MustRegister(operationHistogram)
}
