package ristretto_prometheus

import (
	"github.com/dgraph-io/ristretto"
)

// MetricsProvider is a functor contract for cache metrics
type MetricsProvider func() *ristretto.Metrics

// CacheMetricsProvider is a metrics provider using a [ristretto.Cache]
// instance as source.
type CacheMetricsProvider struct {
	cache *ristretto.Cache
}

// NewCacheMetricsProvider creates a metrics provider using the provided
// cache instance as source.
func NewCacheMetricsProvider(cache *ristretto.Cache) *CacheMetricsProvider {
	return &CacheMetricsProvider{
		cache: cache,
	}
}

// Provide is a [MetricsProvider] implementation. It simply
// returns the Metrics field from the intenal cache value.
func (p *CacheMetricsProvider) Provide() *ristretto.Metrics {
	return p.cache.Metrics
}
