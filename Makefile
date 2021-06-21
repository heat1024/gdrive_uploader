TEST ?= $(shell $(GO) list ./... | grep -v vendor)
VERSION = $(shell cat version)
GOVERSION = $(shell go version | awk '{print $$3}')
REVISION = $(shell git describe --always --exclude '*')
DATE = $(shell date '+%Y%m%d-%H%M%S%Z')
INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
GO ?= GO111MODULE=on go
TAG ?= $(shell git tag | grep v$(VERSION) >/dev/null && echo "v$(VERSION)-$(REVISION)" || echo "v$(VERSION)")

ifeq ("$(shell uname)","Darwin")
GORELEASER ?= GO111MODULE=on goreleaser
else
GORELEASERPATH ?= /usr/local/bin/goreleaser
GORELEASER ?= GO111MODULE=on $(GORELEASERPATH)
endif

default: build

depsdev: ## Installing dependencies for development
	$(GO) get golang.org/x/lint/golint

goreleaser: # install goreleaser
ifeq ("$(shell uname)","Darwin")
	which goreleaser >/dev/null || brew install goreleaser/tap/goreleaser goreleaser
else
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
	install ./bin/goreleaser $(GORELEASERPATH)
	rm -rf ./bin
endif

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -min_confidence 1.1 -set_exit_status $(TEST)

build: depsdev goreleaser ## Build for release
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Building$(RESET)"
	env GOVERSION=$(GOVERSION) env REVISION=$(REVISION) env DATE=$(DATE) $(GORELEASER) build --snapshot --rm-dist

release: goreleaser
	$(GO) mod tidy
	git tag $(TAG)
	git push origin $(TAG)
	env GOVERSION=$(GOVERSION) env REVISION=$(REVISION) env DATE=$(DATE) $(GORELEASER) release --rm-dist

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(RESET) %s\n", $$1, $$2}'

.PHONY: default lint build depsdev goreleaser release help
