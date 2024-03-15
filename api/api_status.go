package api

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const (
	basePathStatus = "/api/status"
)

// DefaultStatusApi is the default implementation of StatusApi interface
type DefaultStatusApi struct {
	client *resty.Client
}

// NewStatusApi permit to get default implementation of StatusApi interface
func NewStatusApi(client *resty.Client) StatusApi {
	return &DefaultStatusApi{
		client: client,
	}
}

func (h DefaultStatusApi) Status() (status map[string]any, err error) {
	resp, err := h.client.R().Get(basePathStatus)
	if err != nil {
		return nil, err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, NewAPIError(resp.StatusCode(), resp.Status())
	}
	kibanaStatus := make(map[string]any)
	err = json.Unmarshal(resp.Body(), &kibanaStatus)
	if err != nil {
		return nil, err
	}
	log.Debug("KibanaStatus: ", kibanaStatus)

	return kibanaStatus, nil
}
