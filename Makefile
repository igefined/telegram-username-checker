PACKAGE = $(shell go list -m)
BINARY_NAME = uchecker
BUILD_DIR = build
VERSION ?= $(shell git describe --exact-match --tags 2> /dev/null)
COMMIT ?= $(shell git rev-parse HEAD)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")
LDFLAGS = -ldflags "-w -X ${PACKAGE}/internal/app.Version=${VERSION} -X ${PACKAGE}/internal/app.BuildDate=${BUILD_DATE} -X ${PACKAGE}/internal/app.Commit=${COMMIT}"
# golang-ci tag
GOLANGCI_TAG:=1.56.0
# Path to the binary
LOCAL_BIN:=$(CURDIR)/bin
# Path to the binary golang-ci
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
# Minimal Golang version
MIN_GO_VERSION = 1.22.0
LINTER_TIMEOUT:=5m

ifneq (,$(wildcard .env))
	include .env
	export
endif

##################### Checks to run golang-ci #####################
# Local bin version check
ifneq ($(wildcard $(GOLANGCI_BIN)),)
GOLANGCI_BIN_VERSION:=$(shell $(GOLANGCI_BIN) --version)
ifneq ($(GOLANGCI_BIN_VERSION),)
GOLANGCI_BIN_VERSION_SHORT:=$(shell echo "$(GOLANGCI_BIN_VERSION)" | sed -E 's/.* version (.*) built from .* on .*/\1/g')
else
GOLANGCI_BIN_VERSION_SHORT:=0
endif
ifneq "$(GOLANGCI_TAG)" "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_BIN_VERSION_SHORT)))"
GOLANGCI_BIN:=
endif
endif

# Global bin version check
ifneq (, $(shell which golangci-lint))
GOLANGCI_VERSION:=$(shell golangci-lint --version 2> /dev/null )
ifneq ($(GOLANGCI_VERSION),)
GOLANGCI_VERSION_SHORT:=$(shell echo "$(GOLANGCI_VERSION)"|sed -E 's/.* version (.*) built from .* on .*/\1/g')
else
GOLANGCI_VERSION_SHORT:=0
endif
ifeq "$(GOLANGCI_TAG)" "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_VERSION_SHORT)))"
GOLANGCI_BIN:=$(shell which golangci-lint)
endif
endif
##################### End of golang-ci checks #####################

# Install linter
.PHONY: install-lint
install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info "Downloading golangci-lint v$(GOLANGCI_TAG)")
	tmp=$$(mktemp -d) && cd $$tmp && pwd && go mod init temp && go get -d github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG) && \
		go build -ldflags "-X 'main.version=$(GOLANGCI_TAG)' -X 'main.commit=test' -X 'main.date=test'" -o $(LOCAL_BIN)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint && \
		rm -rf $$tmp
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint
endif

# Linter will check only diffs with main branch (default)
.PHONY: lint
lint: install-lint
	$(GOLANGCI_BIN) run --config=.golangci.yml ./... --new-from-rev=origin/main --timeout=$(LINTER_TIMEOUT) --build-tags=$(SERVICE_NAME)

# Run full code lint
.PHONY: lint-full
lint-full: lint
	$(GOLANGCI_BIN) run --config=.golangci.yml ./... --build-tags=$(SERVICE_NAME)

# Linter will check only diffs with main branch and auto fix them.
.PHONY: lint-fix
lint-fix: lint
	$(GOLANGCI_BIN) run --fix --config=.golangci.yml ./... --new-from-rev=origin/main --timeout=$(LINTER_TIMEOUT) --build-tags=$(SERVICE_NAME)

# Install config to your home directory.
.PHONY: install-config
install-config:
	@cp .golangci.yml $(HOME)/.golangci.yaml
	@echo "Golangci config installed to $(HOME)/.golangci.yaml"

.PHONY: update
update:
	go mod tidy
	go mod verify

bin/:
	mkdir -p bin

.PHONY: build
build:
	go build -tags '${TAGS}' ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${PACKAGE}/cmd/${BINARY_NAME}

.PHONY: run
run: build
	./bin/uchecker


