package opensearchdashboard

import (
	"crypto/tls"

	"github.com/disaster37/opensearch-dashboard/v2/api"
	"github.com/go-resty/resty/v2"
)

// Client is the client interface to interact with Opensearch Dashboard API
type Client interface {
	api.Api
	Client() *resty.Client
}

// Config contain the value to access on Opensearch Dashboard API
type Config struct {
	Address          string
	Username         string
	Password         string
	DisableVerifySSL bool
	CAs              []string
}

// DefaultClient contain the REST client and the API specification
type DefaultClient struct {
	api.Api
	client *resty.Client
}

// NewDefaultClient init client with empty config
func NewDefaultClient() (Client, error) {
	return NewClient(Config{})
}

// NewClient init client with custom config
func NewClient(cfg Config) (Client, error) {
	if cfg.Address == "" {
		cfg.Address = "http://localhost:5601"
	}

	restyClient := resty.New().
		SetBaseURL(cfg.Address).
		SetBasicAuth(cfg.Username, cfg.Password).
		SetHeader("osd-xsrf", "true").
		SetHeader("securitytenant", "default").
		SetHeader("Content-Type", "application/json")

	for _, path := range cfg.CAs {
		restyClient.SetRootCertificate(path)
	}

	client := &DefaultClient{
		client: restyClient,
		Api:    api.New(restyClient),
	}

	if cfg.DisableVerifySSL {
		client.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return client, nil

}

func (h *DefaultClient) Client() *resty.Client {
	return h.client
}
