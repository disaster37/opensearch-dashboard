package samples

import (
	"encoding/json"
	"log"

	"github.com/disaster37/opensearch-dashboard/v2/api"
	"k8s.io/utils/ptr"
)

func SavedObjectCreate() {
	client := getClient()

	indexPatternObject := &api.Object{
		Type: "index-pattern",
		Id:   "test",
		Attributes: map[string]any{
			"title": "test-pattern-*",
		},
	}

	resp, err := client.SavedObject().Create("", indexPatternObject, true)
	if err != nil {
		log.Fatalf("Error creating object: %s", err)
	}
	log.Println(resp)

}

func SavedObjectGet() {
	client := getClient()

	resp, err := client.SavedObject().Get("", "index-pattern", "test")
	if err != nil {
		log.Fatalf("Error getting index pattern save object: %s", err)
	}
	log.Println(resp)
}

func SavedObjectFind() {
	client := getClient()

	search := api.SavedObjectFindOption{
		Search:       "test",
		SearchFields: []string{"id"},
	}
	resp, err := client.SavedObject().Find("", search)
	if err != nil {
		log.Fatalf("Error searching index pattern: %s", err)
	}
	log.Println(resp)
}

func SavedObjectUpdate() {
	client := getClient()

	indexPatternObject := &api.Object{
		Type: "index-pattern",
		Id:   "test",
		Attributes: map[string]any{
			"title": "test-pattern-bis-*",
		},
	}

	resp, err := client.SavedObject().Update("", indexPatternObject)
	if err != nil {
		log.Fatalf("Error updating index pattern: %s", err)
	}
	log.Println(resp)
}

func SavedObjectExport() {
	client := getClient()

	request := api.SavedObjectExportOption{
		IncludeReferencesDeep: ptr.To[bool](true),
		ExcludeExportDetails:  ptr.To[bool](true),
		Objects: []api.SavedObjectBulkGetOption{
			{
				Type: "index-patter",
				Id:   "test",
			},
		},
	}
	resp, err := client.SavedObject().Export("", request)
	if err != nil {
		log.Fatalf("Error exporting index pattern: %s", err)
	}
	log.Println(resp)
}

func SavedObjectImport() {
	client := getClient()

	indexPatternObject := &api.Object{
		Type: "index-pattern",
		Id:   "test",
		Attributes: map[string]any{
			"title": "test-pattern-*",
		},
	}
	b, err := json.Marshal(indexPatternObject)
	if err != nil {
		panic(err)
	}

	option := api.SavedObjectImportOption{
		Overwrite: true,
	}
	resp, err := client.SavedObject().Import("", option, b)
	if err != nil {
		log.Fatalf("Error importing index pattern: %s", err)
	}
	log.Println(resp)
}

func SavedObjectDelete() {
	client := getClient()

	if err := client.SavedObject().Delete("", "index-pattern", "test", false); err != nil {
		log.Fatalf("Error deleting index pattern: %s", err)
	}
	log.Println("Index pattern successfully deleted")
}
