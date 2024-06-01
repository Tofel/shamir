# Define the default target
.PHONY: test-unit
test_uni: test-go test-python clean

# Run Go tests
.PHONY: test-go
test-go:
	@echo "Running Go tests..."
	@go test ./... -v || { echo "Go tests failed!"; exit 1; }

# Run Python tests
.PHONY: test-python
test-python:
	@echo "Running Python tests..."
	@pytest || { echo "Python tests failed!"; exit 1; }

# Run clean (optional)
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@go clean
	@find . -type f -name "*.pyc" -delete
	@find . -type d -name "__pycache__" -exec rm -r {} +
	@find . -type d -name ".pytest_cache" -exec rm -r {} +

# Define a help target
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make          Run all tests (Go and Python)"
	@echo "  make test     Run all tests (Go and Python)"
	@echo "  make test-go  Run only Go tests"
	@echo "  make test-python Run only Python tests"
	@echo "  make build          Build both Docker images"
	@echo "  make build-go       Build Go Docker image"
	@echo "  make build-python   Build Python Docker image"	
	@echo "  make clean    Clean up generated files"

# Build both Docker images
.PHONY: build
build: build-go build-python

# Build only Go Docker image for multiple platforms
.PHONY: build-go
build-go:
	@echo "Building Go Docker image for linux/amd64 and linux/arm64..."
	@docker buildx create --use
	@docker buildx build --platform linux/amd64 -t shamir-go:latest-amd64 -f Dockerfile-go . --load || { echo "Go Docker build failed!"; exit 1; }
	@docker buildx build --platform linux/arm64 -t shamir-go:latest-arm64 -f Dockerfile-go . --load || { echo "Go Docker build failed!"; exit 1; }

# Build only Python Docker image
.PHONY: build-python
build-python:
	@echo "Building Python Docker image..."
	@docker build -t shamir-python:latest -f Dockerfile-python . || { echo "Python Docker build failed!"; exit 1; }

.PHONY: test-local
test-local:
	@echo "Running local tests..."
	@./test_local.sh

.PHONY: test-docker
test-docker:
	@echo "Running docker tests..."
	@./test_docker.sh

.PHONY: build-test-docker
build-test-docker: build test-docker

.PHONY: test
test: test-unit test-local build-test-docker

