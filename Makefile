
SHELL          = /bin/bash
CFG            = .env

PRG           ?= $(shell basename $$PWD)
REPO          ?= cmd/
# -----------------------------------------------------------------------------
# Build config

GO            ?= go
VERSION       ?= $(shell git describe --tags)
SOURCES       ?= cmd/*/*.go *.go

# -----------------------------------------------------------------------------

.PHONY: all gen doc build-standalone coverage cov-html build test lint fmt vet 

##
## Available make targets
##

# default: show target list
all: help

# ------------------------------------------------------------------------------
## Sources

## Run from sources
run:
	$(GO) run -ldflags "-X main.version=$(VERSION)" ./cmd/$(PRG)/ $(REPO)

vrun:
	$(GO) run -ldflags "-X main.version=$(VERSION)" ./cmd/$(PRG)/ --debug $(REPO)

## Build app with checks
build-all: lint lint-more vet cov build

## Build app
build: 
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/$(PRG)

## Build app used in docker from scratch
build-standalone: cov vet lint lint-more
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=`git describe --tags`" -installsuffix 'static' -a ./cmd/$(PRG)

## Generate mocks
gen:
	$(GO) generate ./...

## Format go sources
fmt:
	$(GO) fmt ./...

## Run vet
vet:
	$(GO) vet ./...

## Run linter
lint:
	golint ./...

## Run more linters
lint-more:
	golangci-lint run ./...

## Run tests and fill coverage.out
cov: cov-clean coverage.out

# internal target
coverage.out: $(SOURCES)
	$(GO) test -test.v -test.race -coverprofile=$@ -covermode=atomic -tags test ./...

## Open coverage report in browser
cov-html: cov
	$(GO) tool cover -html=coverage.out

## Clean coverage report
cov-clean:
	rm -f coverage.*

# ------------------------------------------------------------------------------
## Misc

## Count lines of code (including tests) and update LOC.md
cloc: LOC.md

LOC.md: $(SOURCES)
	cloc --by-file --not-match-f='(_moq_test.go|ml|.md|.sh|.json|file)$$' --md . > $@ 2>/dev/null
	cloc --by-file --not-match-f='(_test.go|ml|.md|.sh|.json|file)$$' . 2>/dev/null
	cloc --by-file --not-match-f='_moq_test.go$$' --match-f='_test.go$$' .  2>/dev/null

## List Makefile targets
help:  Makefile
	@grep -A1 "^##" $< | grep -vE '^--$$' | sed -E '/^##/{N;s/^## (.+)\n(.+):(.*)/\t\2:\1/}' | column -t -s ':'
