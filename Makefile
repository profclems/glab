OS = $(shell uname | tr A-Z a-z)
export PATH := $(abspath bin/):${PATH}


# Build variables
export CGO_ENABLED ?= 0
export GOOS = $(shell go env GOOS)
ifeq (${VERBOSE}, 1)
ifeq ($(filter -v,${GOARGS}),)
	GOARGS += -v
endif
TEST_FORMAT = short-verbose
endif

GLAB_VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
DATE_FMT = +%Y-%m-%d
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif

ifndef CGO_CPPFLAGS
    export CGO_CPPFLAGS := $(CPPFLAGS)
endif
ifndef CGO_CFLAGS
    export CGO_CFLAGS := $(CFLAGS)
endif
ifndef CGO_LDFLAGS
    export CGO_LDFLAGS := $(LDFLAGS)
endif

HASGOTESTSUM := $(shell which gotestsum 2> /dev/null)
HASGOCILINT := $(shell which golangci-lint 2> /dev/null)

ifdef HASGOTESTSUM
    GOTEST=gotestsum
else
    GOTEST=bin/gotestsum
endif

ifdef HASGOCILINT
    GOLINT=golangci-lint
else
    GOLINT=bin/golangci-lint
endif

GO_LDFLAGS := -X main.build=$(BUILD_DATE) $(GO_LDFLAGS)
GO_LDFLAGS := $(GO_LDFLAGS) -X main.version=$(GLAB_VERSION)
GOURL ?= github.com/profclems/glab
BUILDLOC ?= ./bin/glab

# Dependency versions
GOTESTSUM_VERSION = 0.6.0
GOLANGCI_VERSION = 1.32.2

# Add the ability to override some variables
# Use with care
-include override.mk

.PHONY: build
.DEFAULT_GOAL := build
build:
	go build -trimpath -ldflags "$(GO_LDFLAGS) -X main.debugMode=false" -o $(BUILDLOC) $(GOURL)/cmd/glab

clean: ## Clear the working area and the project
	rm -rf ./bin ./.glab-cli ./test/testdata-* ./coverage.txt coverage-*
.PHONY: clean

.PHONY: install
install: ## Install glab in $GOPATH/bin
	GO111MODULE=on go install -trimpath -ldflags "$(GO_LDFLAGS) -X main.debugMode=false" $(GOURL)/cmd/glab

.PHONY: run
run:
	go run -trimpath -ldflags "$(GO_LDFLAGS) -X main.debugMode=true" ./cmd/glab $(run)

.PHONY: rt
rt: ## Test release without publishing
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: rtdebug
rtdebug: ## Test release with debug info
	goreleaser --snapshot --skip-publish --rm-dist --debug

.PHONY: release
release:
	goreleaser $(run)

.PHONY: gen-docs
gen-docs: ## Generate docs
	go run ./cmd/gen-docs/docs.go
	#cp ./docs/glab.rst ./docs/index.rst

.PHONY: check
check: test lint ## Run tests and linters

ifdef HASGOTESTSUM
bin/gotestsum:
	@echo "Skip this"
else
bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum
endif

bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_${OS}_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}

TEST_PKGS ?= ./pkg/... ./internal/... ./commands/... ./cmd/...
.PHONY: test
test: TEST_FORMAT ?= short
test: SHELL = /bin/bash
test: export CGO_ENABLED=1
test:  bin/gotestsum ## Run tests
	$(GOTEST) --no-summary=skipped --junitfile ./coverage.xml --format ${TEST_FORMAT} -- -coverprofile=./coverage.txt -covermode=atomic $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)


ifdef HASGOCILINT
bin/golangci-lint:
	echo "Skip this"
else
bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
endif

bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	$(GOLINT) run

.PHONY: fix
fix: bin/golangci-lint ## Fix lint violations
	go mod tidy
	$(GOLINT) run --fix
	gofmt -s -w .
	goimports -w .

# Add custom targets here
-include custom.mk

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'