PKG_NAME=kbapi
DASHBOARD_URL = http://127.0.0.1:5601
DASHBOARD_USERNAME = admin
DASHBOARD_PASSWORD = vLPeJYa8.3RqtZCcAK6jNz

all: help


test: fmt
	DASHBOARD_URL=${DASHBOARD_URL} DASHBOARD_USERNAME=${DASHBOARD_USERNAME} DASHBOARD_PASSWORD=${DASHBOARD_PASSWORD} go test ./... -v -count 1 -parallel 1 -race -coverprofile=coverage.out -covermode=atomic $(TESTARGS) -timeout 120m

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./