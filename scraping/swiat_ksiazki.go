package scraping

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type RequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type SwiatKsiazkiScraper struct {
	http RequestDoer
}

func NewSwiatKsiazkiScraper() *SwiatKsiazkiScraper {
	return &SwiatKsiazkiScraper{
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *SwiatKsiazkiScraper) Price(query string) (*BookPrice, error) {
	item, err := s.queryItemFromAPI(query)
	if err != nil {
		return nil, err
	}

	return s.mapGraphqlItemToPrice(item), nil
}

func (s *SwiatKsiazkiScraper) buildQueryURL(query string) string {
	q := url.Values{}
	q.Add("hash", "3409055357")
	q.Add("sort_1", `{"position":"ASC"}`)
	q.Add("filter_1", `{"price":{},"category_id":{"eq":"442"},"customer_group_id":{"eq":"0"}}`)
	q.Add("search_1", fmt.Sprintf(`"%s"`, query))
	q.Add("pageSize_1", "1")
	q.Add("currentPage_1", "1")

	scrapeURL, _ := url.Parse("https://www.swiatksiazki.pl/graphql")
	scrapeURL.RawQuery = q.Encode()

	return scrapeURL.String()
}

type swiatKsiazkiGraphQLResponse struct {
	Data struct {
		Products struct {
			TotalCount int                       `json:"total_count"`
			Items      []swiatKsiazkiGraphQLItem `json:"items"`
		} `json:"products"`
	} `json:"data"`
}

type swiatKsiazkiGraphQLItem struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	StockItem struct {
		InStock bool `json:"in_stock"`
	} `json:"stock_item"`
	PriceRange struct {
		MinimumPrice struct {
			FinalPrice struct {
				Currency string  `json:"currency"`
				Value    float32 `json:"value"`
			} `json:"final_price"`
		} `json:"minimum_price"`
	} `json:"price_range"`
	SmallImage struct {
		URL string `json:"url"`
	} `json:"small_image"`
	Dictionary struct {
		Authors []struct {
			Name string `json:"name"`
		} `json:"authors"`
	} `json:"dictionary"`
}

func (s *SwiatKsiazkiScraper) queryItemFromAPI(query string) (*swiatKsiazkiGraphQLItem, error) {
	queryURL := s.buildQueryURL(query)
	log.Println("visiting", queryURL)

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response *swiatKsiazkiGraphQLResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	if response.Data.Products.TotalCount < 1 || len(response.Data.Products.Items) == 0 {
		return nil, ErrNoResult
	}
	log.Println("finished scraping:", queryURL)

	return &response.Data.Products.Items[0], nil
}

func (s *SwiatKsiazkiScraper) mapGraphqlItemToPrice(item *swiatKsiazkiGraphQLItem) *BookPrice {
	author := "-"
	if len(item.Dictionary.Authors) != 0 {
		author = item.Dictionary.Authors[0].Name
	}

	return &BookPrice{
		Title:    item.Name,
		Author:   author,
		Price:    fmt.Sprintf("%.2f zł", item.PriceRange.MinimumPrice.FinalPrice.Value),
		URL:      fmt.Sprintf("https://www.swiatksiazki.pl%s", item.URL),
		CoverURL: item.SmallImage.URL,
	}
}
