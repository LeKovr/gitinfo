## gitinfo Makefile:
## Get git repo metagata (lib) and generate gitinfo.json via go generate (cmd)
#:

SHELL          = /bin/sh
CFG            = .env
PRG           ?= $(shell basename $$PWD)
REPO          ?= cmd/

# -----------------------------------------------------------------------------
# Build config

GO            ?= go
SOURCES       ?= cmd/*/*.go *.go
VERSION       ?= $(shell git describe --tags --always)
GODOC_REPO    ?= github.com/pgmig/$(PRG)
# -----------------------------------------------------------------------------

.PHONY: all gen doc build-standalone coverage cov-html build test lint fmt vet 

# default: show target list
all: help

# ------------------------------------------------------------------------------
## Compile operations
#:

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
	golangci-lint run ./...

## Run tests and fill coverage.out
test: clean coverage.out

# internal target
coverage.out: $(SOURCES)
	$(GO) test -race -coverprofile=$@ -covermode=atomic -tags test ./...

## Open coverage report in browser
cov-html: test
	$(GO) tool cover -html=coverage.out

## Run from sources
run:
	$(GO) run -ldflags "-X main.version=$(VERSION)" ./cmd/$(PRG)/ $(REPO)

vrun:
	$(GO) run -ldflags "-X main.version=$(VERSION)" ./cmd/$(PRG)/ --debug $(REPO)

## Build app with checks
build-all: vet lint test build

## Build app
build: $(SOURCES)
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/$(PRG)

## Clean coverage report
clean:
	rm -f coverage.*

# ------------------------------------------------------------------------------
## Other
#:

## update docs at pkg.go.dev
godoc:
	vf=$(VERSION) ; v=$${vf%%-*} ; echo "Update for $$v..." ; \
	curl 'https://proxy.golang.org/$(GODOC_REPO)/@v/'$$v'.info'

# This code handles group header and target comment with one or two lines only
## list Makefile targets
## (this is default target)
help:
	@grep -A 1 -h "^## " $(MAKEFILE_LIST) \
  | sed -E 's/^--$$// ; /./{H;$$!d} ; x ; s/^\n## ([^\n]+)\n(## (.+)\n)*(.+):(.*)$$/"    " "\4" "\1" "\3"/' \
  | sed -E 's/^"    " "#" "(.+)" "(.*)"$$/"" "" "" ""\n"\1 \2" "" "" ""/' \
  | xargs printf "%s\033[36m%-15s\033[0m %s %s\n"
