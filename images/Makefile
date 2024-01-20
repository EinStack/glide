VENDOR ?= einstack
PROJECT ?= Glide
SOURCE ?= https://github.com/EinStack/glide
LICENSE ?= Apache-2.0
DESCRIPTION ?= "A lightweight, cloud-native LLM gateway"
REPOSITORY ?= einstack/glide

VERSION ?= dev
RC_PART ?= rc
COMMIT ?= $(shell git describe --dirty --always --abbrev=15)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")


# OCI Labels: https://specs.opencontainers.org/image-spec/annotations
# Test images via: docker run --rm --platform linux/amd64 -i einstack/glide:dev-alpine --config config.dev.yaml

alpine: ## Make an alpine-based image
	@echo "🛠️ Build alpine image ($(VERSION)).."
	@echo "- Commit: $(COMMIT)"
	@echo "- Build Date: $(BUILD_DATE)"
	@docker build .. -t $(REPOSITORY):$(VERSION)-alpine -f alpine.Dockerfile \
	--build-arg VERSION="$(VERSION)" \
	--build-arg COMMIT="$(COMMIT)" \
	--build-arg BUILD_DATE="$(BUILD_DATE)" \
	--label=org.opencontainers.image.vendor="$(VENDOR)" \
	--label=org.opencontainers.image.title="$(PROJECT)" \
	--label=org.opencontainers.image.revision="$(COMMIT)" \
	--label=org.opencontainers.image.version="$(VERSION)" \
	--label=org.opencontainers.image.created="$(BUILD_DATE)" \
	--label=org.opencontainers.image.source="$(SOURCE)" \
	--label=org.opencontainers.image.licenses="$(LICENSE)" \
	--label=org.opencontainers.image.description=$(DESCRIPTION)

ubuntu: ## Make an ubuntu-based image
	@echo "🛠️ Build ubuntu image ($(VERSION)).."
	@echo "- Commit: $(COMMIT)"
	@echo "- Build Date: $(BUILD_DATE)"
	@docker build .. -t $(REPOSITORY):$(VERSION)-ubuntu -f ubuntu.Dockerfile \
	--build-arg VERSION="$(VERSION)" \
	--build-arg COMMIT="$(COMMIT)" \
	--build-arg BUILD_DATE="$(BUILD_DATE)" \
	--label=org.opencontainers.image.vendor=$(VENDOR) \
	--label=org.opencontainers.image.title=$(PROJECT) \
	--label=org.opencontainers.image.revision="$(COMMIT)" \
	--label=org.opencontainers.image.version="$(VERSION)" \
	--label=org.opencontainers.image.created="$(BUILD_DATE)" \
	--label=org.opencontainers.image.source=$(SOURCE) \
	--label=org.opencontainers.image.licenses=$(LICENSE) \
	--label=org.opencontainers.image.description=$(DESCRIPTION)

distroless: ## Make an distroless-based image
	@echo "🛠️ Build distroless image ($(VERSION)).."
	@echo "- Commit: $(COMMIT)"
	@echo "- Build Date: $(BUILD_DATE)"
	@docker build .. -t $(REPOSITORY):$(VERSION)-distroless -f distroless.Dockerfile \
	--build-arg VERSION="$(VERSION)" \
	--build-arg COMMIT="$(COMMIT)" \
	--build-arg BUILD_DATE="$(BUILD_DATE)" \
	--label=org.opencontainers.image.vendor=$(VENDOR) \
	--label=org.opencontainers.image.title=$(PROJECT) \
	--label=org.opencontainers.image.revision="$(COMMIT)" \
	--label=org.opencontainers.image.version="$(VERSION)" \
	--label=org.opencontainers.image.created="$(BUILD_DATE)" \
	--label=org.opencontainers.image.source=$(SOURCE) \
	--label=org.opencontainers.image.licenses=$(LICENSE) \
	--label=org.opencontainers.image.description=$(DESCRIPTION)

redhat: ## Make an Red Hat-based image
	@echo "🛠️ Build Red Hat image ($(VERSION)).."
	@echo "- Commit: $(COMMIT)"
	@echo "- Build Date: $(BUILD_DATE)"
	@docker build .. -t $(REPOSITORY):$(VERSION)-redhat -f redhat.Dockerfile \
	--build-arg VERSION="$(VERSION)" \
	--build-arg COMMIT="$(COMMIT)" \
	--build-arg BUILD_DATE="$(BUILD_DATE)" \
	--label=org.opencontainers.image.vendor=$(VENDOR) \
	--label=org.opencontainers.image.title=$(PROJECT) \
	--label=org.opencontainers.image.revision="$(COMMIT)" \
	--label=org.opencontainers.image.version="$(VERSION)" \
	--label=org.opencontainers.image.created="$(BUILD_DATE)" \
	--label=org.opencontainers.image.source=$(SOURCE) \
	--label=org.opencontainers.image.licenses=$(LICENSE) \
	--label=org.opencontainers.image.description=$(DESCRIPTION)

all: alpine ubuntu distroless redhat

publish-ghcr-%: ## Push images to Github Registry
	@echo "🚚Pushing the $* image to Github Registry.."
	@docker tag $(REPOSITORY):$(VERSION)-$* ghcr.io/$(REPOSITORY):$(VERSION)-$*
	@echo "- pushing ghcr.io/$(REPOSITORY):$(VERSION)-$*"
	@docker push ghcr.io/$(REPOSITORY):$(VERSION)-$*
	@echo $(VERSION) | grep -q $(RC_PART) || { \
		docker tag $(REPOSITORY):$(VERSION)-$* ghcr.io/$(REPOSITORY):latest-$*; \
		echo "- pushing ghcr.io/$(REPOSITORY):latest-$*"; \
		docker push ghcr.io/$(REPOSITORY):latest-$*; \
		if [ "$*" = "alpine" ]; then \
			docker tag $(REPOSITORY):$(VERSION)-$* ghcr.io/$(REPOSITORY):latest; \
			echo "- pushing ghcr.io/$(REPOSITORY):latest"; \
			docker push ghcr.io/$(REPOSITORY):latest; \
		fi; \
	}
