package main

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
)

type Cache struct {
	*cache.Cache
	items prometheus.Gauge
}

func (c *Cache) Run() {
	for {
		c.items.Set(float64(c.Cache.ItemCount()))
		time.Sleep(10 * time.Second)
	}
}

func newCache() (*Cache, error) {
	c := &Cache{}
	c.Cache = cache.New(cache.NoExpiration, 10*time.Minute)
	c.items = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_items_cached",
		Help: "Number of items currently in cache",
	})
	prometheus.MustRegister(c.items)
	return c, nil
}
