package scraping

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

type SwiatKsiazkiScraper struct {
}

func NewSwiatKsiazkiScraper() *SwiatKsiazkiScraper {
	return &SwiatKsiazkiScraper{}
}

func (s *SwiatKsiazkiScraper) Price(query string) (*BookPrice, error) {
	c := colly.NewCollector()
	wg := sync.WaitGroup{}
	var result *BookPrice

	c.OnRequest(func(r *colly.Request) {
		wg.Add(1)
		log.Println("visiting", r.URL)
	})

	c.OnError(func(res *colly.Response, err error) {
		wg.Done()
		log.Println(fmt.Sprintf("failed to scrape swiat_ksiazki for '%s':", query), err)
	})

	c.OnHTML(".product-items > .product-item:first-of-type", func(e *colly.HTMLElement) {
		linkEl := e.DOM.Find("a.product-item-link")
		link, _ := linkEl.Attr("href")

		price := e.DOM.Find(".price-box .special-price").Text()
		// Fallback for when a product item isn't on sale and therefore doesn't have a "special-price".
		if price == "" {
			price = e.DOM.Find(".price-box").Text()
		}
		if price == "" {
			price = "N/A"
		}

		result = &BookPrice{
			Title: linkEl.Text(),
			Price: parsePrice(price),
			URL:   link,
		}
	})

	c.OnScraped(func(r *colly.Response) {
		wg.Done()
		log.Println("finished scraping:", r.Request.URL)
	})

	scrapeURL := s.buildPriceScrapingURL(query)
	if err := c.Visit(scrapeURL); err != nil {
		return nil, errors.Wrap(err, "failed to scrape price")
	}

	wg.Wait()

	if result == nil {
		return nil, ErrNoResult
	}
	return result, nil
}

func (s *SwiatKsiazkiScraper) buildPriceScrapingURL(query string) string {
	q := url.Values{}
	q.Add("cat", "4")
	q.Add("q", query)

	scrapeURL, _ := url.Parse("https://www.swiatksiazki.pl/catalogsearch/result")
	scrapeURL.RawQuery = q.Encode()

	return scrapeURL.String()
}
