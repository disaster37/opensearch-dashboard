package samples

import "log"

func Status() {
	client := getClient()

	resp, err := client.Status().Status()
	if err != nil {
		log.Fatalf("Error getting status: %s", err)
	}
	log.Println(resp)
}
