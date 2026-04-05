# Define the default target
.PHONY: test
test: test-unit test-local build-test-docker

# Run unit tests
.PHONY: test-unit
test_unit: test-go test-python clean

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
	@echo "  make          		 		Run all tests (unit, local and docker)"
	@echo "  make test     		 		Run all tests (unit, local and docker)"
	@echo "  make test-go  		 		Run only Go unit tests"
	@echo "  make test-python 	 		Run only Python unit tests"
	@echo "  make test-local 	 		Run local interoperability tests"
	@echo "  make test-docker 	 		Run docker interoperability tests"
	@echo "  make build          		Build both Docker images"
	@echo "  make build-go       		Build Go Docker image"
	@echo "  make build-python  	  	Build Python Docker image"	
	@echo "  make build-test-docker   	Build Docker images and run Docker tests"	
	@echo "  make clean    		 		Clean up generated files"

# Build both Docker images
.PHONY: build
build: build-go build-python

# Build only Go Docker image for multiple platforms
.PHONY: build-go
build-go:
	@echo "Building multiarch Go Docker image..."
	@docker buildx inspect multiarch >/dev/null 2>&1 \
  	&& docker buildx use multiarch \
  	|| docker buildx create --name multiarch --driver docker-container --use
	@docker buildx inspect --bootstrap

	@docker buildx build \
	--platform linux/amd64,linux/arm64 \
	-t shamir-go:latest \
	-f Dockerfile-go \
	--load .	

# Build only Python Docker image
.PHONY: build-python
build-python:
	@echo "Building multiarch Python Docker image..."
	@docker buildx inspect multiarch >/dev/null 2>&1 \
  	&& docker buildx use multiarch \
  	|| docker buildx create --name multiarch --driver docker-container --use
	@docker buildx inspect --bootstrap

	@docker buildx build \
	--platform linux/amd64,linux/arm64 \
	-t shamir-python:latest \
	-f Dockerfile-python \
	--load .		

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
