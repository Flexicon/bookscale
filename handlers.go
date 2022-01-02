package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/flexicon/bookscale/scraping"
	"github.com/labstack/echo/v4"
)

func IndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index", IndexTplArgs{})
}

func SearchHandler(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("q"))
	results := NewSearchResults(query)

	if query == "" {
		return c.Redirect(http.StatusFound, "/")
	}

	wg := sync.WaitGroup{}
	for source, scraper := range scraping.PriceScrapers {
		wg.Add(1)

		go func(source string, scraper scraping.PriceScraper) {
			defer wg.Done()

			price, err := scraper.Price(query)
			if err != nil {
				results.AddError(source, err)
				return
			}
			results.AddPrice(source, price)
		}(source, scraper)
	}

	wg.Wait()

	return c.Render(http.StatusOK, "index", IndexTplArgs{results})
}

type IndexTplArgs struct {
	SearchResults *SearchResults
}

type SearchResults struct {
	Query  string
	Prices map[string]*scraping.BookPrice
	Errors map[string]error

	mu sync.Mutex
}

// NewSearchResults constructor.
func NewSearchResults(query string) *SearchResults {
	return &SearchResults{
		Query:  query,
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
