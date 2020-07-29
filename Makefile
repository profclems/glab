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

GO_LDFLAGS := -X github.com/cli/cli/command.Version=$(GLAB_VERSION) $(GO_LDFLAGS)
GO_LDFLAGS := -X github.com/cli/cli/command.BuildDate=$(BUILD_DATE) $(GO_LDFLAGS)

build:
	go build -trimpath -ldflags "$(GO_LDFLAGS)" -o ./bin/glab ./cmd/glab

run:
	go run -trimpath -ldflags "$(GO_LDFLAGS)" cmd/glab/main.go $(var)

test:
	go test ./...

rt: #Test release
	goreleaser --snapshot --skip-publish --rm-dist

release:
	goreleaser $(var)

compileall:
	mkdir -p ./bin
	mkdir -p ./bin/$(GLAB_VERSION)
	cp cmd/glab/main.go ./bin/$(GLAB_VERSION)/glab.go
	./scripts/compile-all-plaforms.bash ./bin/$(GLAB_VERSION)/glab.go
