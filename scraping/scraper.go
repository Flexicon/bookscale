package scraping

// PriceScrapers map for all available price scrapers.
var PriceScrapers = map[string]PriceScraper{
	"Allegro": NewAllegroScraper(),
}

type PriceScrape struct {
	Price string
	URL   string
}

type PriceScraper interface {
	Price(query string) (*PriceScrape, error)
}
