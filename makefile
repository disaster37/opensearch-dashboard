PKG_NAME=kbapi
DASHBOARD_URL = http://127.0.0.1:5601
DASHBOARD_USERNAME = admin
DASHBOARD_PASSWORD = vLPeJYa8.3RqtZCcAK6jNz

all: help

.PHONY: test
test: fmt
	DASHBOARD_URL=${DASHBOARD_URL} DASHBOARD_USERNAME=${DASHBOARD_USERNAME} DASHBOARD_PASSWORD=${DASHBOARD_PASSWORD} go test ./... -v -count 1 -parallel 1 -race -coverprofile=coverage.out -covermode=atomic $(TESTARGS) -timeout 120m

.PHONY: fmt
fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./

.PHONY: mock-gen
mock-gen:
	go install go.uber.org/mock/mockgen@v0.3.0
	mockgen --build_flags=--mod=mod -destination=mocks/client.go -package=mocks github.com/disaster37/opensearch-dashboard/v2 Client
	mockgen --build_flags=--mod=mod -destination=mocks/api.go -package=mocks github.com/disaster37/opensearch-dashboard/v2/api Api,SavedObjectApi,ShortenUrlApi,StatusApi
