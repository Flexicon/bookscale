package main

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

// Cache wrapper.
type Cache struct {
	cache *cache.Cache
}

// Get retrieves an item from cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// Set adds an item to the cache, replacing any existing items for the given key.
func (c *Cache) Set(key string, val interface{}) {
	c.cache.Set(key, val, cache.DefaultExpiration)
}

// priceCache instance.
var priceCache *Cache

// InitCache sets up the app cache.
func InitCache() error {
	cacheTTL := time.Duration(viper.GetInt64("cache.ttl")) * time.Second
	priceCache = &Cache{cache.New(cacheTTL, 10*time.Minute)}
	return nil
}
