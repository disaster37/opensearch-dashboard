package samples

import (
	"log"

	opensearchdashboard "github.com/disaster37/opensearch-dashboard/v2"
)

func getClient() opensearchdashboard.Client {
	cfg := opensearchdashboard.Config{
		Address:          "http://127.0.0.1:5601",
		Username:         "admin",
		Password:         "admin",
		DisableVerifySSL: true,
	}

	client, err := opensearchdashboard.NewClient(cfg)

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
		panic(err)
	}

	status, err := client.Status().Status()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	log.Println(status)

	return client
}
