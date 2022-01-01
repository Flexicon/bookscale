package main

import (
	"github.com/flexicon/bookscale/scraping"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
)

func IndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

type SearchResults struct {
	SearchTerm string
	Prices     map[string]*scraping.BookPrice
	Errors     []error

	mu sync.Mutex
}

// NewSearchResults constructor.
func NewSearchResults(searchTerm string) *SearchResults {
	return &SearchResults{
		SearchTerm: searchTerm,
		Prices:     make(map[string]*scraping.BookPrice),
		Errors:     make([]error, 0),
	}
}

func (r *SearchResults) AddPrice(source string, price *scraping.BookPrice) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Prices[source] = price
}

func (r *SearchResults) AddError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Errors = append(r.Errors, err)
}

func SearchHandler(c echo.Context) error {
	searchTerm := c.QueryParam("term")
	results := NewSearchResults(searchTerm)

	wg := sync.WaitGroup{}
	for source, scraper := range scraping.PriceScrapers {
		wg.Add(1)

		go func(source string, scraper scraping.PriceScraper) {
			defer wg.Done()

			price, err := scraper.Price(searchTerm)
			if err != nil {
				results.AddError(err)
				return
			}
			results.AddPrice(source, price)
		}(source, scraper)
	}

	wg.Wait()

	return c.Render(http.StatusOK, "search", results)
}
