package ristretto_prometheus_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"

	ristretto_prometheus "github.com/wolfmetr/ristretto-prometheus"
)

func TestNewCollector(t *testing.T) {
	cache, err := ristretto.NewCache[string, string](&ristretto.Config[string, string]{
		NumCounters: 1e3,
		MaxCost:     1 << 30,
		BufferItems: 64,
		Metrics:     true, // enable ristretto metrics
	})
	if err != nil {
		panic(err)
	}
	defer cache.Close()

	ristrettoCollector, err := ristretto_prometheus.NewMetricsCollector(
		cache.Metrics,
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

	// fill the cache
	for i := 0; i < 123; i++ {
		cache.Set(
			fmt.Sprintf("key%d", i),
			fmt.Sprintf("val%d", i),
			1,
		)
	}

	// wait for value to pass through buffers
	cache.Wait()

	// generate hits
	for i := 50; i < 99; i++ {
		key := fmt.Sprintf("key%d", i)
		if _, ok := cache.Get(key); !ok {
			t.Errorf("expected key: %s", key)
		}
	}

	// generate misses
	for i := 150; i < 170; i++ {
		key := fmt.Sprintf("key%d", i)
		if _, ok := cache.Get(key); ok {
			t.Errorf("unexpected key: %s", key)
		}
	}

	expected := `
		# HELP appname_subsystemname_ristretto_hits_ratio The percentage of successful Get calls (hits).
		# TYPE appname_subsystemname_ristretto_hits_ratio gauge
		appname_subsystemname_ristretto_hits_ratio{app_version="v1.2.3"} 0.7101449275362319
		# HELP appname_subsystemname_ristretto_hits_total The number of Get calls where a value was found for the corresponding key.
		# TYPE appname_subsystemname_ristretto_hits_total counter
		appname_subsystemname_ristretto_hits_total{app_version="v1.2.3"} 49
		# HELP appname_subsystemname_ristretto_keys_added_total The number of added keys in the cache.
		# TYPE appname_subsystemname_ristretto_keys_added_total counter
		appname_subsystemname_ristretto_keys_added_total{app_version="v1.2.3"} 123
		# HELP appname_subsystemname_ristretto_misses_total The number of Get calls where a value was not found for the corresponding key.
		# TYPE appname_subsystemname_ristretto_misses_total counter
		appname_subsystemname_ristretto_misses_total{app_version="v1.2.3"} 20
	`

	if err := testutil.CollectAndCompare(ristrettoCollector, strings.NewReader(expected)); err != nil {
		t.Errorf("unexpected collecting result:\n%s", err)
	}
}
