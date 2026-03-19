.PHONY: build lint test clean install release snapshot pre-commit pre-commit-install

# Build the binary
build:
	go build -o menlo ./cmd/menlo

# Install binary locally
install: build
	install -m 755 menlo /usr/local/bin/

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
	rm -f menlo
	rm -rf dist/

# Install locally
install: build
	install -m 755 menlo /usr/local/bin/

# Install lefthook hooks
pre-commit-install:
	lefthook install

# Run pre-commit checks
pre-commit:
	lefthook run pre-commit

# Create a release
release:
	goreleaser release --clean

# Create a snapshot release
snapshot:
	goreleaser build --snapshot --clean