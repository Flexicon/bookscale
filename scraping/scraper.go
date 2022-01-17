package scraping

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

const (
	// NoCoverURL used for when no cover was found for a scraped book.
	NoCoverURL = "/img/no_cover.svg"
)

// ErrNoResult is returned by Scrapers if no data is found during scraping.
var ErrNoResult = errors.New("no result")

// PriceScrapers map for all available price scrapers.
var PriceScrapers = map[string]PriceScraper{
	//"allegro": NewAllegroScraper(), // TODO: re-enable after implementing work around properly
	"swiat_ksiazki": NewSwiatKsiazkiScraper(),
	"empik":         NewEmpikScraper(),
}

// BookPrice represents a scraped price item.
type BookPrice struct {
	Title    string
	Price    string
	URL      string
	CoverURL string
}

// PriceScraper allows scraping prices by query.
type PriceScraper interface {
	Price(query string) (*BookPrice, error)
}

// Sources returns a sorted slice of all available scraping source keys.
func Sources() []string {
	var keys []string
	for key, _ := range PriceScrapers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

// parsePrice reads and formats a given price in PLN (zł) from the given input string.
//
// Return format: "24,99 zł"
func parsePrice(input string) string {
	input = strings.TrimSpace(input)
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return ""
	}

	return fmt.Sprintf("%s zł", strings.TrimSpace(parts[0]))
}
