BUILD_PATH ?= _build
PACKAGES_PATH ?= _packages

VERSION := $(shell git describe --tags --always)
ARCH := $(shell go env GOARCH)
OS := $(shell go env GOOS)
PACKAGE_NAME := kuiper-crrc-$(VERSION)-$(OS)-$(ARCH)
GO              := GO111MODULE=on go

TARGET ?= lfedge/ekuiper

export KUIPER_SOURCE := $(shell pwd)

.PHONY: build
build: build_without_edgex

.PHONY: build_prepare
build_prepare:
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/bin
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/etc
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/etc/sources
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/etc/sinks
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/etc/services
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/etc/templates
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/etc/templates/function
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/etc/services/schemas
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/data
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/plugins
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/plugins/sources
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/plugins/sinks
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/plugins/functions
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/plugins/portable
	@mkdir -p $(BUILD_PATH)/$(PACKAGE_NAME)/log

	@cp -r etc/* $(BUILD_PATH)/$(PACKAGE_NAME)/etc
	@cp -r plugins/* $(BUILD_PATH)/$(PACKAGE_NAME)/plugins

.PHONY: build_without_edgex
build_without_edgex: build_prepare
	GO111MODULE=on CGO_ENABLED=1 go build -ldflags="-s -w -X main.Version=$(VERSION) -X main.LoadFileType=relative" -o kuiper cmd/kuiper/main.go
	GO111MODULE=on CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.Version=$(VERSION) -X main.LoadFileType=relative" -o kuiperd cmd/kuiperd/main.go
	@if [ ! -z $$(which upx) ]; then upx ./kuiper; upx ./kuiperd; fi
	@mv ./kuiper ./kuiperd $(BUILD_PATH)/$(PACKAGE_NAME)/bin
	@echo "Build successfully"

.PHONY: pkg_without_edgex
pkg_without_edgex: build_without_edgex
	@make real_pkg

.PHONY: build_with_edgex
build_with_edgex: build_prepare
	GO111MODULE=on CGO_ENABLED=1 go build -ldflags="-s -w -X main.Version=$(VERSION) -X main.LoadFileType=relative" -tags "edgex include_nats_messaging" -o kuiper cmd/kuiper/main.go
	GO111MODULE=on CGO_ENABLED=1 go build -trimpath -ldflags="-s -w -X main.Version=$(VERSION) -X main.LoadFileType=relative" -tags "edgex include_nats_messaging" -o kuiperd cmd/kuiperd/main.go
	@if [ ! -z $$(which upx) ]; then upx ./kuiper; upx ./kuiperd; fi
	@mv ./kuiper ./kuiperd $(BUILD_PATH)/$(PACKAGE_NAME)/bin
	@echo "Build successfully"

.PHONY: build_with_edgex_and_script
build_with_edgex_and_script: build_prepare
	GO111MODULE=on CGO_ENABLED=1 go build -ldflags="-s -w -X main.Version=$(VERSION) -X main.LoadFileType=relative" -tags "edgex include_nats_messaging" -o kuiper cmd/kuiper/main.go
	GO111MODULE=on CGO_ENABLED=1 go build -trimpath -ldflags="-s -w -X main.Version=$(VERSION) -X main.LoadFileType=relative" -tags "edgex include_nats_messaging script" -o kuiperd cmd/kuiperd/main.go
	@if [ ! -z $$(which upx) ]; then upx ./kuiper; upx ./kuiperd; fi
	@mv ./kuiper ./kuiperd $(BUILD_PATH)/$(PACKAGE_NAME)/bin
	@echo "Build successfully"

.PHONY: pkg_with_edgex
pkg_with_edgex: build_with_edgex
	@make real_pkg

.PHONY: build_core
build_core: build_prepare
	GO111MODULE=on CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.Version=$(VERSION) -X main.LoadFileType=relative" -tags core -o kuiperd cmd/kuiperd/main.go
	@if [ ! -z $$(which upx) ]; then upx ./kuiperd; fi
	@mv ./kuiperd $(BUILD_PATH)/$(PACKAGE_NAME)/bin
	@echo "Build successfully"

.PHONY: pkg_core
pkg_core: build_core
	@mkdir -p $(PACKAGES_PATH)
	@cd $(BUILD_PATH) && zip -rq $(PACKAGE_NAME)-core.zip $(PACKAGE_NAME)
	@cd $(BUILD_PATH) && tar -czf $(PACKAGE_NAME)-core.tar.gz $(PACKAGE_NAME)
	@mv $(BUILD_PATH)/$(PACKAGE_NAME)-core.zip $(BUILD_PATH)/$(PACKAGE_NAME)-core.tar.gz $(PACKAGES_PATH)
	@echo "Package core success"

.PHONY: real_pkg
real_pkg:
	@mkdir -p $(PACKAGES_PATH)
	@cd $(BUILD_PATH) && zip -rq $(PACKAGE_NAME).zip $(PACKAGE_NAME)
	@cd $(BUILD_PATH) && tar -czf $(PACKAGE_NAME).tar.gz $(PACKAGE_NAME)
	@mv $(BUILD_PATH)/$(PACKAGE_NAME).zip $(BUILD_PATH)/$(PACKAGE_NAME).tar.gz $(PACKAGES_PATH)
	@echo "Package build success"

.PHONY: docker
docker:
	@docker build -t $(PACKAGE_NAME):latest -f deploy/docker/Dockerfile-slim-python .
	@docker save -o $(PACKAGES_PATH)/$(PACKAGE_NAME)-docker.tar $(PACKAGE_NAME):latest

.PHONY: clean
clean:
	@rm -rf cross_build.tar linux_amd64 linux_arm64 linux_arm_v7 linux_386
	@rm -rf _build _packages _plugins

tidy:
	@echo "go mod tidy"
	go mod tidy && git diff go.mod go.sum

lint:tools/lint/bin/golangci-lint
	@echo "linting"
	tools/lint/bin/golangci-lint run ./... ./extensions/... ./tools/kubernetes/...
	cd sdk/go && ../../tools/lint/bin/golangci-lint run

tools/lint/bin/golangci-lint:
	GOBIN=$(shell pwd)/tools/lint/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
