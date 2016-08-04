VERSION ?= $(shell cat VERSION)
GIT_HASH := $(shell git rev-parse HEAD)

XC_OS ?= darwin linux
XC_ARCH ?= 386 amd64
LDFLAGS :=-X main.Version=$(VERSION)

all: test

deps:
	@go get -t ./...

clean:
	@rm -rf dist/*

build:
	@gox -os="$(XC_OS)" -arch="$(XC_ARCH)" -ldflags "$(LDFLAGS)" -output "dist/{{.OS}}_{{.Arch}}/{{.Dir}}"

release: build
	@for platform in $$(find ./dist -mindepth 1 -maxdepth 1 -type d); do \
		pushd $$platform >/dev/null; \
		zip ../$$(basename $$platform).zip ./* >/dev/null; \
		popd >/dev/null; \
	done
	@ghr $(VERSION) dist/

test: deps
	@echo "# Running tests"
	@go test ./...

.PHONY: clean release test

