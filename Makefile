.PHONY: build clean test install release

BINARY_NAME=gh-projects
BINARY_PATH=./cmd/gh-projects
VERSION ?= $(shell git describe --tags --always --dirty)

build:
	go build -ldflags="-X main.version=$(VERSION)" -o $(BINARY_NAME) $(BINARY_PATH)

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test ./...

install:
	go install $(BINARY_PATH)

run: build
	./$(BINARY_NAME)

release:
	@echo "To create a release:"
	@echo "1. Create and push a tag: git tag v1.0.0 && git push origin v1.0.0"
	@echo "2. GitHub Actions will automatically build and create the release"

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  install  - Install the binary to GOPATH/bin"
	@echo "  run      - Build and run the binary"
	@echo "  release  - Show instructions for creating a release"
	@echo "  help     - Show this help message"