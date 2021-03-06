package helmclient

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	PrometheusNamespace = "helmclient"
	PrometheusSubsystem = "library"
)

var (
	errorGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "error_total",
			Help:      "Number of helmclient errors.",
		},
		[]string{"event"},
	)
	eventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "event_total",
			Help:      "Number of helmclient events.",
		},
		[]string{"event", "release"},
	)
	histogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "event",
			Help:      "Histogram for events within the helmclient library.",
		},
		[]string{"event"},
	)
)

func init() {
	prometheus.MustRegister(errorGauge)
	prometheus.MustRegister(eventCounter)
	prometheus.MustRegister(histogram)
}
