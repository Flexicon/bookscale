package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/flexicon/bookscale/scraping"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
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
	query := strings.ToLower(strings.TrimSpace(c.QueryParam("q")))
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
			cachedPrice, hit := priceCache.Get(cacheKey)
			if hit {
				log.Printf("cache hit: %+v", cacheKey)
				results.AddCached(source, cachedPrice)
				return
			}

			price, err := scraper.Price(query)
			if err != nil {
				// Only cache ErrNoResult to avoid scraper spam
				if errors.Is(err, scraping.ErrNoResult) {
					priceCache.Set(cacheKey, err)
				}
				results.AddError(source, err)
				return
			}

			priceCache.Set(cacheKey, price)
			results.AddPrice(source, price)
		}(source, scraper)
	}

	wg.Wait()

	return c.Render(http.StatusOK, "index", IndexTplArgs{
		Query:         query,
		Sources:       scraping.Sources(),
		SearchResults: results,
		NoCoverURL:    scraping.NoCoverURL(),
		BaseImageURL:  viper.GetString("static_asset_base_url"),
	})
}

// IndexTplArgs represents the arguments that are passed to the index template.
type IndexTplArgs struct {
	Sources       []string
	Query         string
	SearchResults *SearchResults
	NoCoverURL    string
	BaseImageURL  string
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

// AddCached value to result concurrently.
func (r *SearchResults) AddCached(source string, val interface{}) {
	switch v := val.(type) {
	case *scraping.BookPrice:
		r.AddPrice(source, v)
	case error:
		r.AddError(source, v)
	default:
		r.AddError(source, scraping.ErrNoResult)
	}
}
