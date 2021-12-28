package scraping

type AllegroScraper struct {
}

func NewAllegroScraper() *AllegroScraper {
	return &AllegroScraper{}
}

func (s *AllegroScraper) Price(query string) (*PriceScrape, error) {
	return &PriceScrape{
		Price: "29,99 zł",
		URL:   "http://go-colly.org/docs/introduction/start/",
	}, nil
}
