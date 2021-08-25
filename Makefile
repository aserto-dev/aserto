SHELL 	   := $(shell which bash)

NO_COLOR   :=\033[0m
OK_COLOR   :=\033[32;01m
ERR_COLOR  :=\033[31;01m
WARN_COLOR :=\033[36;01m
ATTN_COLOR :=\033[33;01m

VERSION    ?= $(shell svu)
COMMIT     ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DATE       ?= $(shell date "+%FT%T%z")

GOARCH     ?= amd64
GOOS       := $(shell go env GOOS)
CGO_ENABLED:=0
LDBASE     := github.com/aserto-dev/aserto/pkg/version
LDFLAGS    := -ldflags "-X ${LDBASE}.ver=${VERSION} -X ${LDBASE}.date=${DATE} -X ${LDBASE}.commit=${COMMIT}"
GOPATH     := $(shell go env GOPATH)

TARGET     := aserto
ROOT_DIR   ?= $(shell git rev-parse --show-toplevel)
BIN_DIR    := ${ROOT_DIR}/bin
SRC_DIR    := ${ROOT_DIR}/cmd
DIST_DIR   := ${ROOT_DIR}/dist
BIN_FILE   := ${BIN_DIR}/${GOOS}-${GOARCH}/${TARGET}$(if $(findstring ${GOOS},windows),".exe","")

${BIN_DIR}:
	@echo -e "${ATTN_COLOR}==> create BIN_DIR ${BIN_DIR} ${NO_COLOR}"
	@mkdir -p ${BIN_DIR}

TESTER     := ${GOPATH}/bin/gotestsum
${TESTER}:
	@echo -e "${ATTN_COLOR}==> $@ ${NO_COLOR}"
	@go install gotest.tools/gotestsum@v1.7.0

LINTER	   := ${GOPATH}/bin/golangci-lint
${LINTER}:
	@echo -e "${ATTN_COLOR}==> $@  ${NO_COLOR}"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1

RELEASER   := ${GOPATH}/bin/goreleaser
${RELEASER}:
	@echo -e "${ATTN_COLOR}==> $@  ${NO_COLOR}"
	@go install github.com/goreleaser/goreleaser@v0.174.2

.PHONY: all
all: deps build test lint

.PHONY: deps
deps:
	@echo -e "${ATTN_COLOR}==> $@ ${NO_COLOR}"
	@go install github.com/caarlos0/svu@v1.7.0

.PHONY: build
build: ${RELEASER}
	@echo -e "${ATTN_COLOR}==> $@ ENV=${ENVIRONMENT} GOOS=${GOOS} GOARCH=${GOARCH} VERSION=${VERSION} COMMIT=${COMMIT} DATE=${DATE} ${NO_COLOR}"
	@${RELEASER} build --snapshot --rm-dist --single-target --config .goreleaser-prod.yml

.PHONY: test 
test: ${TESTER}
	@echo -e "${ATTN_COLOR}==> $@ ${NO_COLOR}"
	@${TESTER} --format short-verbose -- -coverprofile=cover.out -coverpkg=./... -count=1 -timeout 90s -v ${ROOT_DIR}/...

.PHONY: lint
lint: ${LINTER}
	@echo -e "${ATTN_COLOR}==> $@ ${NO_COLOR}"
	@${LINTER} run
	@echo -e "${NO_COLOR}\c"

.PHONY: release
release: ${RELEASER}
	@echo -e "${ATTN_COLOR}==> $@ ${NO_COLOR}"
	@${RELEASER} release --skip-publish --rm-dist --snapshot --config .goreleaser-prod.yml

.PHONY: publish
publish: ${RELEASER}
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN environment variable is undefined)
endif
ifndef ASERTO_TAP
	$(error ASERTO_TAP environment variable is undefined)
endif
	@echo -e "${ATTN_COLOR}==> $@ ${NO_COLOR}"
	@${RELEASER} release --config .goreleaser-publish.yml --rm-dist

.PHONY: clean
clean:
	@echo -e "${ATTN_COLOR}==> $@ ${NO_COLOR}"
	@rm -rf ${BIN_DIR}
	@rm -rf $(DIST_DIR)
