# Variables
NAME:= mts-rss-aggregator
IMG ?= $(NAME):latest
IMG_NAME := $(shell echo $(IMG) | awk -F: '{print $$1}')

##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Development

.PHONY: run
run: ## Run the application with cli arguments: config.yml secrets 8080
	go run main.go config.yaml

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run the tests
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html


.PHONY: build
build: fmt vet ## Build the application for the native architecture
	 CGO_ENABLED=0 go build -v -trimpath -tags netgo -ldflags="-s -w" -o ${NAME} main.go


.PHONY: docker-build
docker-build: build ## Build the Docker image
	docker build -t ${IMG} .