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

func WithKeysAddedMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_keys_added_total",
		Help:      "The number of added keys in the cache.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.KeysAdded()) },
	})
}

func WithKeysUpdatedMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_keys_updated_total",
		Help:      "The number of updated keys in the cache.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.KeysUpdated()) },
	})
}

func WithKeysEvictedMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_keys_evicted_total",
		Help:      "The number of evicted keys from the cache.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.KeysEvicted()) },
	})
}

func WithCostAddedMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_cost_added_total",
		Help:      "The sum of costs that have been added.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.CostAdded()) },
	})
}

func WithCostEvictedMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_cost_evicted_total",
		Help:      "The sum of all costs that have been evicted.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.CostEvicted()) },
	})
}

func WithSetsDroppedMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_sets_dropped_total",
		Help:      "The number of Set calls that don't make it into internal buffers.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.SetsDropped()) },
	})
}

func WithSetsRejectedMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_sets_rejected_total",
		Help:      "The number of Set calls rejected by the policy.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.SetsRejected()) },
	})
}

func WithGetsDroppedMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_gets_dropped_total",
		Help:      "The number of Get counter increments that are dropped.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.GetsDropped()) },
	})
}

func WithGetsKeptMetric() Option {
	return WithMetric(Desc{
		Name:      "ristretto_gets_kept_total",
		Help:      "The number of Get counter increments that are kept.",
		ValueType: prometheus.CounterValue,
		Extractor: func(m *ristretto.Metrics) float64 { return float64(m.GetsKept()) },
	})
}

func WithMetric(d Desc) Option {
	return func(o *config) {
		o.metrics = append(o.metrics, d)
	}
}
