package main

import (
	"time"

	"github.com/patrickmn/go-cache"
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
	priceCache = &Cache{cache.New(15*time.Minute, 10*time.Minute)}
	return nil
}
