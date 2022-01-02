package scraping

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

type EmpikScraper struct {
}

func NewEmpikScraper() *EmpikScraper {
	return &EmpikScraper{}
}

func (s *EmpikScraper) Price(query string) (*BookPrice, error) {
	c := colly.NewCollector()
	wg := sync.WaitGroup{}
	var result *BookPrice

	c.OnRequest(func(r *colly.Request) {
		wg.Add(1)
		log.Println("visiting", r.URL)
	})

	c.OnError(func(res *colly.Response, err error) {
		wg.Done()
		log.Println(fmt.Sprintf("failed to scrape empik for '%s':", query), err)
	})

	c.OnHTML(".search-list-item:first-of-type", func(e *colly.HTMLElement) {
		linkEl := e.DOM.Find(".name .product-title > a")

		price := parsePrice(e.DOM.Find(".price").Text())
		if price == "" {
			price = "N/A"
		}

		link, _ := linkEl.Attr("href")
		// Handle relative URLs by attaching the Empik domain onto the link.
		if strings.Index(link, s.domain()) != 0 {
			link = s.domain() + link
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

	scrapeURL := s.domain() + "/ksiazki,31,s?qtype=basicForm&q=" + query
	if err := c.Visit(scrapeURL); err != nil {
		return nil, errors.Wrap(err, "failed to scrape price")
	}

	wg.Wait()

	if result == nil {
		return nil, ErrNoResult
	}
	return result, nil
}

func (s *EmpikScraper) domain() string {
	return "https://www.empik.com"
}
