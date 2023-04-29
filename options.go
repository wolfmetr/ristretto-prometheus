package ristretto_prometheus

import (
	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricValueExtractor func(m *ristretto.Metrics) float64

// Desc describes metric and function for extract value for it
type Desc struct {
	Name      string
	Help      string
	ValueType prometheus.ValueType

	// Extractor is a function that can extract metric from ristretto.Metrics
	Extractor MetricValueExtractor
}

type config struct {
	namespace string
	subsystem string

	constLabels prometheus.Labels

	metrics []Desc
}

type Option func(*config)

func (c *config) apply(opts []Option) {
	for _, o := range opts {
		o(c)
	}
}

func WithNamespace(namespace string) Option {
	return func(o *config) {
		o.namespace = namespace
	}
}

func WithSubsystem(subsystem string) Option {
	return func(o *config) {
		o.subsystem = subsystem
	}
}

func WithConstLabels(constLabels prometheus.Labels) Option {
	return func(o *config) {
		o.constLabels = constLabels
	}
}

func WithHitsCounterMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_hits_total",
		Help:      "The number of Get calls where a value was found for the corresponding key.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.Hits()) },
	})
}

func WithMissesCounterMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_misses_total",
		Help:      "The number of Get calls where a value was not found for the corresponding key.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.Misses()) },
	})
}

func WithHitsRatioGaugeMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_hits_ratio",
		Help:      "The percentage of successful Get calls (hits).",
		ValueType: prometheus.GaugeValue,
		Extractor: func(m *ristretto.Metrics) float64 { return m.Ratio() },
	})
}

func WithMetric(d Desc) Option {
	return func(o *config) {
		o.metrics = append(o.metrics, d)
	}
}
