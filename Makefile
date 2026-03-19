.PHONY: build lint test clean install release snapshot

# Build the binary
build:
	go build -o menlo-cli ./cmd/menlo-cli

# Build for all platforms
build-all:
	goreleaser build --snapshot --clean

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -f menlo-cli
	rm -rf dist/

# Install locally
install: build
	install -m 755 menlo-cli /usr/local/bin/

# Create a release
release:
	goreleaser release --clean

# Create a snapshot release
snapshot:
	goreleaser build --snapshot --clean