package scraping

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"
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

		link := linkEl.AttrOr("href", "")
		// Handle relative URLs by attaching the Empik domain onto the link.
		if strings.Index(link, s.domain()) != 0 {
			link = s.domain() + link
		}

		coverEl := e.DOM.Find("a.img > img")
		coverURL := coverEl.AttrOr("lazy-img", NoCoverURL)

		result = &BookPrice{
			Title:    linkEl.Text(),
			Price:    price,
			URL:      link,
			CoverURL: coverURL,
		}
	})

	c.OnScraped(func(r *colly.Response) {
		wg.Done()
		log.Println("finished scraping:", r.Request.URL)
	})

	_ = c.Visit(s.buildPriceScrapingURL(query))
	wg.Wait()

	if result == nil {
		return nil, ErrNoResult
	}
	return result, nil
}

func (s *EmpikScraper) domain() string {
	return "https://www.empik.com"
}

func (s *EmpikScraper) buildPriceScrapingURL(query string) string {
	q := url.Values{}
	q.Add("qtype", "basicForm")
	q.Add("q", query)

	scrapeURL, _ := url.Parse(fmt.Sprintf("%s/ksiazki,31,s", s.domain()))
	scrapeURL.RawQuery = q.Encode()

	return scrapeURL.String()
}
