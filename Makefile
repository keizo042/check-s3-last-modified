.PHONU: help build
VERSION := 0.0.1

.DEFAULT_GOAL := build
help: ## help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


build: ## build binary
	@go build -ldflags "-X main.Version=${VERSION}"
