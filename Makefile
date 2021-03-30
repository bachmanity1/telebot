PACKAGE = telebot
BUILDPATH ?= $(CURDIR)
DATE    ?= $(shell date +%FT%T%z)
PKGS	= $(or $(PKG),$(shell cd $(BASE) && cat pkg.list))
BASE	= $(BUILDPATH)
BIN		= $(BASE)/bin
DESTDIR	= /opt/$(PACKAGE)
DOCKERPORT = 10619

#GOENV   = CGO_LDFLAGS_ALLOW='-fopenmp'
UNAME := $(shell uname)
ifeq ($(UNAME), Linux)
	GOENV   ?= CGO_ENABLED=0 GOOS=linux
endif
GOBUILD = ${GOENV} go
GO      = go
GODOC   = godoc
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")
GOV = $(shell $(GO) version)
SUDO = ""
ifneq ($(shell id -u -r),0)
SUDO = sudo -E
endif

BUILDTAG=-tags 'release'
export GO111MODULE=on

.PHONY: all
all: test staticcheck fmt vet lint build ; $(info $(M) building all steps… ) @ ## Build all steps

.PHONY: check
check: test staticcheck fmt vet lint ; $(info $(M) check all steps… ) @ ## Check all steps

.PHONY: build
build: vendor ; $(info $(M) building executable… ) @ ## Build program binary
	$Q cd $(BASE) && $(GOBUILD) build -i \
		$(BUILDTAG) \
		-mod vendor \
		-o $(BIN)/$(PACKAGE)

.PHONY: docker
docker: ; $(info $(M) building docker image… ) @ ## Build for docker image 
	$Q $(SUDO) docker build --cache-from cp/$(PACKAGE):latest --build-arg BUILD_PORT=$(DOCKERPORT) -t cp/$(PACKAGE) .

# Tools

$(BIN):
	@mkdir -p $@

$(BIN)/%: $(BIN) ; $(info $(M) downloading $(REPOSITORY)…)
	$Q cd $(BASE); \
		GOBIN=$(BIN) $(GO) get $(REPOSITORY) \
		|| ret=$$?; \
		exit $$ret

GOLINT = $(BIN)/golint
$(BIN)/golint: REPOSITORY=golang.org/x/lint/golint

GOFMT = $(BIN)/goreturns
$(BIN)/goreturns: REPOSITORY=github.com/sqs/goreturns

GOCOVMERGE = $(BIN)/gocovmerge
$(BIN)/gocovmerge: REPOSITORY=github.com/wadey/gocovmerge

GOCOV = $(BIN)/gocov
$(BIN)/gocov: REPOSITORY=github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: REPOSITORY=github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go2xunit
$(BIN)/go2xunit: REPOSITORY=github.com/tebeka/go2xunit

GOCHECK = $(BIN)/staticcheck
$(BIN)/staticcheck: REPOSITORY=honnef.co/go/tools/cmd/staticcheck

PROTOC = `which protoc`

# Tests

.PHONY: test
test: ; $(info $(M) running go test…) @ ## Run go test on all source files
	$Q cd $(BASE) && ret=0 && for d in $(PKGS); do \
		cd $$d ; \
		$(GO) test -mod vendor -failfast -race -cover $(BUILDTAG) || ret=$$? ; \
		[ "$$ret" != "0" ] && exit $$ret; \
		cd .. ; \
	 done ; exit $$ret

.PHONY: lint
lint: $(GOLINT) ; $(info $(M) running go lint…) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret

.PHONY: vet
vet: ; $(info $(M) running go vet…) @ ## Run go vet on all source files
	$Q cd $(BASE) && ret=0 && for d in $(PKGS); do \
		cd $$d ; \
		$(GO) vet -mod vendor $(BUILDTAG) || ret=$$? ; \
		cd .. ; \
	 done ; exit $$ret

.PHONY: fmt
fmt: $(GOFMT) ; $(info $(M) running go fmt…) @ ## Run gofmt on all source files
	$Q cd $(BASE) && ret=0 && for d in $(PKGS); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

.PHONY: staticcheck
staticcheck: $(GOCHECK) ; $(info $(M) running go check…) @ ## Run gocheck on all source files
	$Q cd $(BASE) && $(GOCHECK) `cat pkg.list`


# Dependency management
#

.PHONY: proto
proto:
ifeq (, $(shell which $(PROTOC)))
	$(info $(M) no protoc compiler… install protoc)
else
	$(info $(M) gen file from proto)
	$Q cd $(BASE)/idl && \
	$(PROTOC) \
		-I. \
	 	--go_out=plugins=grpc:$(BASE)/model \
	 	--go_opt=paths=source_relative *.proto 
endif

.PHONY: pkg.list
pkg.list: 
ifeq (,$(wildcard $(CURDIR)/pkg.list))
	$(info $(M) generate pkg.list file…)
	$Q cd $(BASE) && GO111MODULE=on $(GO) list -f '{{.Dir}}' ./... 2>&1 | grep -v "^go: " | grep -v "^$(PACKAGE)/vendor/" | grep -v nocompile | grep -v logs | grep -v testset | grep -v idl | grep -v model > pkg.list
endif

.PHONY: vendor-check
vendor-check: ; $(info $(M) go.mod file…)
ifneq (,$(wildcard $(CURDIR)/go.mod))
	$(info $(M) no needs mod init…)
else
	$(info $(M) mod init…)
	$Q cd $(BASE) && $(GO) mod init $(PACKAGE)
endif

vendor: pkg.list vendor-check go.mod; $(info $(M) retrieving dependencies…)
	$Q cd $(BASE) && $(GO) mod tidy && $(GO) mod vendor


# Misc

.PHONY: install
install: ; $(info $(M) installing…)	@
	@mkdir -p $(DESTDIR)
	@cp -fp bin/$(PACKAGE) $(DESTDIR)/.
	@cp -fp config.json $(DESTDIR)/config.json

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin
	@rm -rf cmd/docs
	@rm -rf vendor
	@rm -rf pkg.list
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

