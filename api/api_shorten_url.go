package api

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const (
	basePathShortenURL = "/api/shorten_url" // Base URL to access on Kibana shorten URL
)

// ShortenUrlApiCreatePayload is the payload to send when create shorten URL
type ShortenUrlApiCreatePayload struct {
	Url string `json:"url"`
}

// ShortenUrlApiCreateResponse is the response when create shorten URL
type ShortenUrlApiCreateResponse struct {
	UrlId string `json:"urlId"`
}

// DefaultShortenUrlApi is the default implementation of ShortenUrlApi interface
type DefaultShortenUrlApi struct {
	client *resty.Client
}

// NewShortenUrlApi permit to get default implementation of ShortenUrlApi interface
func NewShortenUrlApi(client *resty.Client) ShortenUrlApi {
	return &DefaultShortenUrlApi{
		client: client,
	}
}

func (h DefaultShortenUrlApi) Create(url string) (shortUrl string, err error) {
	if url == "" {
		return "", NewAPIError(600, "You must provide shorten URL object")
	}
	log.Debugf("Shorten URL: %s", url)

	resp, err := h.client.R().
		SetBody(ShortenUrlApiCreatePayload{Url: url}).
		Post(basePathShortenURL)
	if err != nil {
		return "", err
	}

	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		return "", NewAPIError(resp.StatusCode(), resp.Status())
	}

	shortenURLResponse := &ShortenUrlApiCreateResponse{}
	err = json.Unmarshal(resp.Body(), shortenURLResponse)
	if err != nil {
		return "", err
	}
	log.Debug("ShortenURLResponse: ", shortenURLResponse.UrlId)

	return shortenURLResponse.UrlId, nil
}
