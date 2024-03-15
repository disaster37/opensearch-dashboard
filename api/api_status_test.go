package api

import "github.com/stretchr/testify/assert"

func (s *ApiTestSuite) TestStatus() {

	// Create new shorten URL
	status, err := s.Api.Status().Status()
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), status)
}
