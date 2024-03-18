package opensearchdashboard

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type DashboardTestSuite struct {
	suite.Suite
}

func (s *DashboardTestSuite) SetupSuite() {

	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

}

func TestDashboardTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardTestSuite))
}

func (s *DashboardTestSuite) TestNewClient() {

	cfg := Config{
		Address:          "http://127.0.0.1:5601",
		Username:         "admin",
		Password:         "vLPeJYa8.3RqtZCcAK6jNz",
		DisableVerifySSL: true,
	}

	client, err := NewClient(cfg)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), client)

}

func (s *DashboardTestSuite) TestNewDefaultClient() {

	client, err := NewDefaultClient()

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), client)

}
