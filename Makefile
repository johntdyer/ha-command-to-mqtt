.PHONY: build run clean test docker-build docker-run

# Build the application
build:
	go build -o ha-command-to-mqtt .

# Run the application
run:
	go run main.go

# Run with race detection
run-race:
	go run -race main.go

# Clean build artifacts
clean:
	rm -f ha-command-to-mqtt

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build Docker image
docker-build:
	docker build -t ha-command-to-mqtt:latest .

# Run Docker container
docker-run:
	docker run --rm -v $(PWD)/config.yaml:/root/config.yaml ha-command-to-mqtt:latest

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o ha-command-to-mqtt-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o ha-command-to-mqtt-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o ha-command-to-mqtt-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o ha-command-to-mqtt-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o ha-command-to-mqtt-windows-amd64.exe .