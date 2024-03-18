[![build](https://github.com/disaster37/opensearch-dashboard/actions/workflows/workflow.yml/badge.svg)](https://github.com/disaster37/opensearch-dashboard/actions/workflows/workflow.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/disaster37/opensearch-dashboard/v2)](https://goreportcard.com/report/github.com/disaster37/opensearch-dashboard/v2)
[![GoDoc](https://godoc.org/github.com/disaster37/opensearch-dashboard/v2?status.svg)](http://godoc.org/github.com/disaster37/opensearch-dashboard/v2)
[![codecov](https://codecov.io/gh/disaster37/opensearch-dashboard/graph/badge.svg?token=S2EVN8N79U)](https://codecov.io/gh/disaster37/opensearch-dashboard)


# opensearch-dashboard

Golang library to consume Opensearch dashboard Api

It support the following APIs:
  - Saved objects APIs
  - Shorten URL


## Compatibility

It work with the following Opensearch dashboard verison:
  - 2.x

## Installation

Get librairy with gomod:
```bash
go get -u https://github.com/disaster37/opensearch-dashboard
```

## Usage

### Init the client

```go
cfg := opensearchdashboard.Config{
    Address:          "http://127.0.0.1:5601",
    Username:         "admin",
    Password:         "admin",
    DisableVerifySSL: true,
}

client, err := opensearchdashboard.NewClient(cfg)

if err != nil {
    log.Fatalf("Error creating the client: %s", err)
}

status, err := client.Status().Status()
if err != nil {
    log.Fatalf("Error getting response: %s", err)
}
log.Println(status)
```

### Handle shorten URL

```go
// Shorten long URL
shortenURL, err := client.ShortenUrl().Create( "/app/dashboards#/dashboard?_g=()&_a=(description:'',filters:!(),fullScreenMode:!f,options:(hidePanelTitles:!f,useMargins:!t),panels:!((embeddableConfig:(),gridData:(h:15,i:'1',w:24,x:0,y:0),id:'8f4d0c00-4c86-11e8-b3d7-01146121b73d',panelIndex:'1',type:visualization,version:'7.0.0-alpha1')),query:(language:lucene,query:''),timeRestore:!f,title:'New%20Dashboard',viewMode:edit)")
if err != nil {
    log.Fatalf("Error creating shorten URL: %s", err)
}
log.Println(fmt.Sprintf("http://localhost:5601/goto/%s", shortenURL))
```






### Handle save object

```go
// Create new index pattern in default user space
indexPatternObject := &api.Object{
  Type: "index-pattern",
  Id: "test",
  Attributes: map[string]any{
    "title": "test-pattern-*",
  },
}

resp, err := client.SavedObject().Create("", indexPatternObject, true)
if err != nil {
  log.Fatalf("Error creating object: %s", err)
}
log.Println(resp)

// Get index pattern save object from default user space
resp, err := client.SavedObject().Get("", "index-pattern", "test")
if err != nil {
  log.Fatalf("Error getting index pattern save object: %s", err)
}
log.Println(resp)

// Search index pattern from default user space
search := api.SavedObjectFindOption{
  Search: "test",
  SearchFields: []string{"id"},
}
resp, err := client.SavedObject().Find("", search)
if err != nil {
  log.Fatalf("Error searching index pattern: %s", err)
}
log.Println(resp)

// Update index pattern in default user space
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

// Export index pattern from default user space
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

// import index pattern in default user space
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

// Delete index pattern in default user space
if err := client.SavedObject().Delete("", "index-pattern", "test", false); err != nil {
  log.Fatalf("Error deleting index pattern: %s", err)
}
log.Println("Index pattern successfully deleted")
```

### Handle status

```go
resp, err := client.Status().Status()
if err != nil {
  log.Fatalf("Error getting status: %s", err)
}
log.Println(resp)
```

## Contribute

PR are always welcome here !

Please start from the branch `2.x`, implement code and don't forget to implement test to prove it's work as expected.

To test:
```bash
docker-compose up -d
make test
```