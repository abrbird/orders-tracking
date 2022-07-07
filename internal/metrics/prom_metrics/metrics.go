package prom_metrics

import (
	cnfg "github.com/abrbird/orders-tracking/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	Errors      prometheus.Counter
	KafkaErrors prometheus.Counter
}

func New(cfg *cnfg.Config) *Metrics {
	mtrcs := &Metrics{
		Errors: promauto.NewCounter(
			prometheus.CounterOpts{
				Name:        "errors",
				Help:        "Number of common errors.",
				ConstLabels: prometheus.Labels{"service": cfg.Application.Name},
			},
		),
		KafkaErrors: promauto.NewCounter(
			prometheus.CounterOpts{
				Name:        "kafka_errors",
				Help:        "Number of Kafka errors.",
				ConstLabels: prometheus.Labels{"service": cfg.Application.Name},
			},
		),
	}
	return mtrcs
}

func (m Metrics) Error() {
	m.Errors.Inc()
}

func (m Metrics) KafkaError() {
	m.KafkaErrors.Inc()
}
