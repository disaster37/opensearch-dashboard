package api

import (
	"github.com/stretchr/testify/assert"
)

func (s *ApiTestSuite) TestShortenURL() {

	// Create new shorten URL
	shortenURL, err := s.Api.ShortenUrl().Create("/app/kibana#/dashboard?_g=()&_a=(description:'',filters:!(),fullScreenMode:!f,options:(hidePanelTitles:!f,useMargins:!t),panels:!((embeddableConfig:(),gridData:(h:15,i:'1',w:24,x:0,y:0),id:'8f4d0c00-4c86-11e8-b3d7-01146121b73d',panelIndex:'1',type:visualization,version:'7.0.0-alpha1')),query:(language:lucene,query:''),timeRestore:!f,title:'New%20Dashboard',viewMode:edit)")
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), shortenURL)
}
