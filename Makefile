GLAB_VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)

build:
	go build -o bin/glab cmd/glab/main.go

run:
	go run cmd/glab/main.go $(var)

push:
	git push origin $(git branch | sed -n -e 's/^\* \(.*\)/\1/p')

compileall:
	mkdir ./bin && mkdir ./bin/$(GLAB_VERSION)
	cp cmd/glab/main.go ./bin/$(GLAB_VERSION)/glab
	./scripts/compile-all-plaforms.bash ./bin/$(GLAB_VERSION)/glab
