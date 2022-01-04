package main

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// PriceCache instance.
var PriceCache *cache.Cache

// InitCache sets up the app cache.
func InitCache() error {
	PriceCache = cache.New(15*time.Minute, 10*time.Minute)

	return nil
}
