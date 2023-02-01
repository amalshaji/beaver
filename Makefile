.PHONY: build-server build-client run-test-server

build-server:
	go build -ldflags="-s -w" -o beaver_server ./cmd/beaver_server

build-client:
	go build -ldflags="-s -w" -o beaver ./cmd/beaver_client

run-test-server:
	go run ./examples/test_api/main.go