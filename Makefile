VERSION_PACKAGE := glide/pkg
COMMIT ?= $(shell git describe --dirty --long --always --abbrev=15)
VERSION ?= "latest" # TODO: pull/pass a real version

LDFLAGS_COMMON := "-s -w -X $(VERSION_PACKAGE).commitSha=$(COMMIT) -X $(VERSION_PACKAGE).version=$(VERSION)"

.PHONY: help

help:
	@echo "🛠️ Glide Dev Commands:\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Lint the source code
	@echo "🧹 Vetting go.mod.."
	@go vet ./...
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 GoCI Lint.."
	@golangci-lint run ./...

run: ## Run Glide
	@go run -ldflags $(LDFLAGS_COMMON) main.go

build: ## Build Glide
	@go build -ldflags $(LDFLAGS_COMMON) -o ./dist/glide
