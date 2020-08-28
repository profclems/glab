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

GO_LDFLAGS := -X main.build=$(BUILD_DATE) $(GO_LDFLAGS)
GO_LDFLAGS := $(GO_LDFLAGS) -X main.version=$(GLAB_VERSION)
GOURL ?= glab
BUILDLOC ?= ./bin/glab

build:
	go build -trimpath -ldflags "$(GO_LDFLAGS) -X  main.usageMode=prod" -o $(BUILDLOC) $(GOURL)/cmd/glab

install:
	GO111MODULE=on go install -trimpath -ldflags "$(GO_LDFLAGS)-sources -X  main.usageMode=prod" $(GOURL)/cmd/glab

run:
	go run -trimpath -ldflags "$(GO_LDFLAGS) -X main.usageMode=dev" ./cmd/glab $(var)

tests:
	bash -c "trap 'trap - SIGINT SIGTERM ERR; rm coverage-* 2>&1 > /dev/null; exit 1' SIGINT SIGTERM ERR; $(MAKE) internal-test"

internal-test:
	rm coverage-* 2>&1 > /dev/null || true
	GO111MODULE=on go test -coverprofile=coverage-main.out -covermode=count -coverpkg ./... -run=$(run)$(GOURL)/cmd/glab $(GOURL)/commands $(GOURL)/internal/...
	go get -u github.com/wadey/gocovmerge
	gocovmerge coverage-*.out > coverage.txt && rm coverage-*.out

rt: #Test release
	goreleaser --snapshot --skip-publish --rm-dist

rtdebug: #Test release
	goreleaser --snapshot --skip-publish --rm-dist --debug

release:
	goreleaser $(var)

gen-docs:
	go run ./cmd/gen-docs/docs.go
	cp ./docs/glab.md ./docs/index.md
