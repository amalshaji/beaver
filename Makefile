.PHONY: build-server build-client run-test-server

build-server:
	go build -ldflags="-s -w" -o beaver_server ./cmd/beaver_server

build-server-image:
	docker buildx build --platform linux/amd64,linux/arm64 -t amalshaji/beaver:latest -f deployments/Dockerfile --push .

build-client:
	go build -ldflags="-s -w" -o beaver ./cmd/beaver_client

run-test-server:
	go run ./examples/test_api/main.go

