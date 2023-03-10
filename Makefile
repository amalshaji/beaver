.PHONY: build-server build-client run-test-server

build-server:
	go build -ldflags="-s -w" -o beaver_server ./cmd/beaver_server

publish-server-image:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t amalshaji/beaver:$(tag) \
		-f deployments/Dockerfile --push .

build-client:
	go build -ldflags="-s -w" -o beaver ./cmd/beaver_client

start-test-servers:
	@echo "Starting test server"
	@go run tests/server.go &
	@sleep 5	
	@echo "Starting beaver server"
	@go run cmd/beaver_server/main.go --config docs/beaver_server.yaml &
	@sleep 5
	@echo "Starting beaver client"
	@go run cmd/beaver_client/main.go --config docs/beaver_client.yaml http 9999 --subdomain test &
	@sleep 5

kill-test-servers:
	@echo "Killing test server"
	@lsof -t -i:9999 | xargs kill -9
	@echo "Killing beaver server"
	@lsof -t -i:8080 | xargs kill -9
