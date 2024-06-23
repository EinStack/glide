CHECKER_BIN=$(PWD)/tmp/bin
VERSION_PACKAGE := glide/pkg/version
COMMIT ?= $(shell git describe --dirty --long --always --abbrev=15)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= "latest"

LDFLAGS_COMMON := "-X $(VERSION_PACKAGE).commitSha=$(COMMIT) -X $(VERSION_PACKAGE).Version=$(VERSION) -X $(VERSION_PACKAGE).buildDate=$(BUILD_DATE)"

.PHONY: help

help:
	@echo "🛠️ Glide Dev Commands:\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


install-checkers: ## Install static checkers
	@echo "🚚 Downloading binaries.."
	@GOBIN=$(CHECKER_BIN) go install mvdan.cc/gofumpt@latest
	@GOBIN=$(CHECKER_BIN) go install golang.org/x/vuln/cmd/govulncheck@latest
	@GOBIN=$(CHECKER_BIN) go install github.com/securego/gosec/v2/cmd/gosec@latest
	@GOBIN=$(CHECKER_BIN) go install github.com/swaggo/swag/cmd/swag@latest

lint: install-checkers ## Lint the source code
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 Formatting files.."
	@go fmt ./...
	@$(CHECKER_BIN)/gofumpt -l -w .
	@echo "🧹 Vetting go.mod.."
	@go vet ./...
	@echo "🧹 GoCI Lint.."
	@golangci-lint run ./...

vuln: install-checkers ## Check for vulnerabilities
	@echo "🔍 Checking for vulnerabilities"
	@#$(CHECKER_BIN)/govulncheck -test ./... enable in https://github.com/EinStack/glide/issues/169
	@$(CHECKER_BIN)/gosec -quiet -exclude=G104 ./...

run: ## Run Glide
	@air -c .air.toml

build: ## Build Glide
	@echo "🔨Building Glide binary.."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@go build -ldflags $(LDFLAGS_COMMON) -o ./dist/glide;

test: ## Run tests
	@go test -v -count=1 -race -shuffle=on -coverprofile=coverage.txt ./...

docs-api: install-checkers ## Generate OpenAPI API docs
	@$(CHECKER_BIN)/swag init

telemetry-up: ## Start observability services needed to receive Glides signals
	@docker-compose --profile telemetry up --wait
	@echo "Jaeger UI: http://localhost:16686/"
	@echo "Grafana UI: http://localhost:3000/"

telemetry-down: ## Shutdown observability services needed to receive Glides signals
	@docker-compose --profile telemetry down
