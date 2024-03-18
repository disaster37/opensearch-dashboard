package samples

import (
	"log"
)

func ShortenUrl() {

	client := getClient()

	shortenURL, err := client.ShortenUrl().Create("/app/dashboards#/dashboard?_g=()&_a=(description:'',filters:!(),fullScreenMode:!f,options:(hidePanelTitles:!f,useMargins:!t),panels:!((embeddableConfig:(),gridData:(h:15,i:'1',w:24,x:0,y:0),id:'8f4d0c00-4c86-11e8-b3d7-01146121b73d',panelIndex:'1',type:visualization,version:'7.0.0-alpha1')),query:(language:lucene,query:''),timeRestore:!f,title:'New%20Dashboard',viewMode:edit)")
	if err != nil {
		log.Fatalf("Error creating shorten URL: %s", err)
	}
	log.Printf("http://localhost:5601/goto/%s", shortenURL)
}
