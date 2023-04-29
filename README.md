# ristretto-prometheus
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/wolfmetr/ristretto-prometheus)

Prometheus Collector for Ristretto Cache metrics

## Usage

### Example
```go
package main

import (
	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"

	ristretto_prometheus "github.com/wolfmetr/ristretto-prometheus"
)

func main() {
    cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e3,
		MaxCost:     1 << 30,
		BufferItems: 64,
		Metrics:     true, // enable ristretto metrics
	})
	if err != nil {
		panic(err)
	}
	defer cache.Close()

	ristrettoCollector, err := ristretto_prometheus.NewCollector(
		cache,
		ristretto_prometheus.WithNamespace("appname"),
		ristretto_prometheus.WithSubsystem("subsystemname"),
		ristretto_prometheus.WithConstLabels(prometheus.Labels{"app_version": "v1.2.3"}),
		ristretto_prometheus.WithHitsCounterMetric(),
		ristretto_prometheus.WithMissesCounterMetric(),
		ristretto_prometheus.WithMetric(ristretto_prometheus.Desc{
			Name:      "ristretto_keys_added_total",
			Help:      "The number of added keys in the cache.",
			ValueType: prometheus.CounterValue,
			Extractor: func(m *ristretto.Metrics) float64 { return float64(m.KeysAdded()) },
		}),
		ristretto_prometheus.WithHitsRatioGaugeMetric(),
	)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(ristrettoCollector)
}

```