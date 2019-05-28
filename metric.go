package helmclient

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	PrometheusNamespace = "operatorkit"
	PrometheusSubsystem = "controller"
)

var (
	controllerErrorGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "error_total",
			Help:      "Number of helmclient errors.",
		},
		[]string{"event", "release"},
	)
	controllerHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "event",
			Help:      "Histogram for events within the helmclient controller.",
		},
		[]string{"event", "release"},
	)
)

func init() {
	prometheus.MustRegister(controllerErrorGauge)
	prometheus.MustRegister(controllerHistogram)
}
