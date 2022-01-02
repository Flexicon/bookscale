package scraping

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type AllegroScraper struct {
}

func NewAllegroScraper() *AllegroScraper {
	return &AllegroScraper{}
}

func (s *AllegroScraper) Price(query string) (*BookPrice, error) {
	token, err := s.getOAuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate oauth token")
	}

	return &BookPrice{
		Price: "35,00 zł",
		URL:   token,
		//URL:   "http://go-colly.org/docs/introduction/start/",
	}, nil
}

func (s *AllegroScraper) getOAuthToken() (string, error) {
	oauthURL := "https://allegro.pl/auth/oauth/token?grant_type=client_credentials"
	req, err := http.NewRequest(http.MethodPost, oauthURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to prepare oauth request")
	}

	req.Header.Add("Authorization", s.encodeClientCredentials())

	res, err := s.httpClient().Do(req)
	if err != nil {
		return "", errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("response status code %d", res.StatusCode))
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", errors.Wrap(err, "failed to parse response body")
	}

	return result.AccessToken, nil
}

func (s *AllegroScraper) encodeClientCredentials() string {
	clientID := viper.GetString("allegro.client_id")
	clientSecret := viper.GetString("allegro.client_secret")

	credentials := fmt.Sprintf("%s:%s", clientID, clientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))

	return "Basic " + encoded
}

func (s *AllegroScraper) httpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
