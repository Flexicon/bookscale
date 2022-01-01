package scraping

import (
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"log"
	"sync"
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
		log.Println("something went wrong:", err)
		log.Printf("%s", res.Body)
	})

	c.OnHTML(".product-items > .product-item:first-of-type", func(e *colly.HTMLElement) {
		linkEl := e.DOM.Find("a.product-item-link")
		link, _ := linkEl.Attr("href")

		price := e.DOM.Find(".price-box .special-price").Text()
		if price == "" {
			price = e.DOM.Find(".price-box").Text()
		}

		result = &BookPrice{
			Title: linkEl.Text(),
			Price: price,
			URL:   link,
		}
	})

	c.OnScraped(func(r *colly.Response) {
		wg.Done()
		log.Println("finished scraping:", r.Request.URL)
	})

	scrapeURL := "https://www.swiatksiazki.pl/catalogsearch/result?cat=4&q=" + query
	if err := c.Visit(scrapeURL); err != nil {
		return nil, errors.Wrap(err, "failed to scrape price")
	}

	wg.Wait()

	return result, nil
}
