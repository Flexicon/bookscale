package scraping

// PriceScrapers map for all available price scrapers.
var PriceScrapers = map[string]PriceScraper{
	//"allegro": NewAllegroScraper(),
	"swiat_ksiazki": NewSwiatKsiazkiScraper(),
}

// BookPrice represents a scraped price item.
type BookPrice struct {
	Title string
	Price string
	URL   string
}

// PriceScraper allows scraping prices by query.
type PriceScraper interface {
	Price(query string) (*BookPrice, error)
}
