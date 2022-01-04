package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/flexicon/bookscale/scraping"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

// SetupRoutes for the app.
func SetupRoutes(e *echo.Echo) {
	e.GET("/", IndexHandler)
	e.GET("/search", SearchHandler)

	e.Static("", "./static")
}

// IndexHandler route handler.
func IndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index", IndexTplArgs{})
}

// SearchHandler route handler.
func SearchHandler(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("q"))
	results := NewSearchResults()

	if query == "" {
		return c.Redirect(http.StatusFound, "/")
	}

	wg := sync.WaitGroup{}
	for source, scraper := range scraping.PriceScrapers {
		wg.Add(1)

		go func(source string, scraper scraping.PriceScraper) {
			defer wg.Done()

			cacheKey := fmt.Sprintf("%s - %s", source, query)
			cachedPrice, hit := PriceCache.Get(cacheKey)
			if hit {
				log.Printf("cache hit: %+v", cacheKey)
				results.AddPrice(source, cachedPrice.(*scraping.BookPrice))
				return
			}

			price, err := scraper.Price(query)
			if err != nil {
				results.AddError(source, err)
				return
			}

			PriceCache.Set(cacheKey, price, cache.DefaultExpiration)
			results.AddPrice(source, price)
		}(source, scraper)
	}

	wg.Wait()

	return c.Render(http.StatusOK, "index", IndexTplArgs{
		Query:         query,
		Sources:       scraping.Sources(),
		SearchResults: results,
	})
}

// IndexTplArgs represents the arguments that are passed to the index template.
type IndexTplArgs struct {
	Sources       []string
	Query         string
	SearchResults *SearchResults
}

// SearchResults holds scraping results and handles adding them concurrently.
type SearchResults struct {
	Prices map[string]*scraping.BookPrice
	Errors map[string]error

	mu sync.Mutex
}

// NewSearchResults constructor.
func NewSearchResults() *SearchResults {
	return &SearchResults{
		Prices: make(map[string]*scraping.BookPrice),
		Errors: make(map[string]error),
	}
}

// AddPrice to results concurrently.
func (r *SearchResults) AddPrice(source string, price *scraping.BookPrice) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Prices[source] = price
}

// AddError to results concurrently.
func (r *SearchResults) AddError(source string, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Errors[source] = err
}
