package scraping

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type TaniaKsiazkaScraper struct {
}

func NewTaniaKsiazkaScraper() *TaniaKsiazkaScraper {
	return &TaniaKsiazkaScraper{}
}

func (s *TaniaKsiazkaScraper) Price(query string) (*BookPrice, error) {
	c := colly.NewCollector()
	wg := sync.WaitGroup{}
	var result *BookPrice

	c.OnRequest(func(r *colly.Request) {
		wg.Add(1)
		log.Println("visiting", r.URL)
	})

	c.OnError(func(res *colly.Response, err error) {
		wg.Done()
		log.Println(fmt.Sprintf("failed to scrape tania_ksiazka for '%s':", query), err)
	})

	c.OnHTML("article ul > li:first-of-type:not(.active)", func(e *colly.HTMLElement) {
		linkEl := e.DOM.Find("a.product-title")

		price := parsePrice(e.DOM.Find(".product-price").Text())
		if price == "" {
			price = "N/A"
		}

		link := linkEl.AttrOr("href", "")
		// Handle relative URLs by attaching the Empik domain onto the link.
		if strings.Index(link, s.domain()) != 0 {
			link = s.domain() + link
		}

		coverEl := e.DOM.Find(".product-image img")
		coverURL := coverEl.AttrOr("data-src", NoCoverURL())

		author := strings.TrimSpace(e.DOM.Find(".product-authors").Text())
		if author == "" {
			author = "-"
		}

		result = &BookPrice{
			Title:    linkEl.Text(),
			Author:   author,
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

func (s *TaniaKsiazkaScraper) domain() string {
	return "https://www.taniaksiazka.pl"
}

func (s *TaniaKsiazkaScraper) buildPriceScrapingURL(query string) string {
	return fmt.Sprintf("%s/Szukaj/q-%s?params[tg]=1&params[last]=tg", s.domain(), url.PathEscape(query))
}
