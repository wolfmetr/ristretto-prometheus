package ristretto_prometheus

import (
	"errors"

	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
)

var ErrDuplicateMetricName = errors.New("duplicate metric name")

// Collector implements the prometheus.Collector interface.
var _ prometheus.Collector = (*Collector)(nil)

type Collector struct {
	provider MetricsProvider

	// metrics contains all descriptions to be registered on a
	// Prometheus metrics registry for the Ristretto cache.
	metrics []metric
}

type metric struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	extractor MetricValueExtractor
}

// NewCollector returns a Prometheus metrics collector using metrics from the
// provided cache instance.
func NewCollector(cache *ristretto.Cache, opts ...Option) (*Collector, error) {
	provider := NewCacheMetricsProvider(cache)

	return NewMetricsCollector(provider.Provide, opts...)
}

// NewCollector returns a Prometheus metrics collector using metrics from the
// given provider.
func NewMetricsCollector(provider MetricsProvider, opts ...Option) (*Collector, error) {
	var conf config
	conf.apply(opts)

	uniqFQNames := make(map[string]struct{})
	metrics := make([]metric, 0, len(conf.metrics))
	for _, c := range conf.metrics {
		fqName := prometheus.BuildFQName(conf.namespace, conf.subsystem, c.Name)
		if _, ok := uniqFQNames[fqName]; ok {
			return nil, ErrDuplicateMetricName
		}
		uniqFQNames[fqName] = struct{}{}

		metrics = append(metrics, metric{
			desc:      prometheus.NewDesc(fqName, c.Help, nil, conf.constLabels),
			valueType: c.ValueType,
			extractor: c.Extractor,
		})
	}

	return &Collector{
		provider: provider,
		metrics:  metrics,
	}, nil
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m.desc
	}
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	metrics := c.provider()
	for _, m := range c.metrics {
		ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.extractor(metrics))
	}
}
