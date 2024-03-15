package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/disaster37/opensearch/v2"
	"github.com/disaster37/opensearch/v2/config"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"k8s.io/utils/ptr"
)

type ApiTestSuite struct {
	suite.Suite
	client *resty.Client
	Api
}

func (s *ApiTestSuite) SetupSuite() {

	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	address := os.Getenv("DASHBOARD_URL")
	username := os.Getenv("DASHBOARD_USERNAME")
	password := os.Getenv("DASHBOARD_PASSWORD")

	if address == "" {
		panic("You need to put opensearch dashboard url on environment variable DASHBOARD_URL. If you need auth, you can use DASHBOARD_USERNAME and DASHBOARD_PASSWORD")
	}

	restyClient := resty.New().
		SetBaseURL(address).
		SetBasicAuth(username, password).
		SetHeader("osd-xsrf", "true").
		SetHeader("Content-Type", "application/json").
		SetHeader("securitytenant", "default").
		SetDebug(false)

	s.client = restyClient
	s.Api = New(restyClient)

	// Wait dashboard is online
	isOnline := false
	nbTry := 0
	for isOnline == false {
		_, err := s.Api.Status().Status()
		if err == nil {
			isOnline = true
		} else {
			logrus.Error(err.Error())
			time.Sleep(5 * time.Second)
			if nbTry == 10 {
				panic(fmt.Sprintf("We wait 50s that Opensearch dashboard start: %s", err))
			}
			nbTry++
		}
	}

	// Create a tenant for test purpose
	cfg := &config.Config{
		URLs:        []string{"https://127.0.0.1:9200"},
		Username:    username,
		Password:    password,
		Sniff:       ptr.To[bool](false),
		Healthcheck: ptr.To[bool](false),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := opensearch.NewClientFromConfig(cfg)
	if err != nil {
		panic(err)
	}
	if _, err := es.SecurityPutTenant("test").Body(opensearch.SecurityPutTenant{Description: ptr.To[string]("test")}).Do(context.Background()); err != nil {
		panic(err)
	}

}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
