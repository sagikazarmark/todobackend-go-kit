# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

OS = $(shell uname | tr A-Z a-z)
export PATH := $(abspath bin/):${PATH}

# Build variables
BUILD_DIR ?= build
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DATE_FMT = +%FT%T%z
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
LDFLAGS += -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}
export CGO_ENABLED ?= 0
ifeq (${VERBOSE}, 1)
ifeq ($(filter -v,${GOARGS}),)
	GOARGS += -v
endif
TEST_FORMAT = short-verbose
endif

# Dependency versions
GOTESTSUM_VERSION = 0.4.2
GOLANGCI_VERSION = 1.27.0
MGA_VERSION = 0.2.0
OPENAPI_GENERATOR_VERSION = 4.3.1
PROTOC_VERSION = 3.12.2
BUF_VERSION = 0.15.0
PROTOC_GEN_KIT_VERSION = 0.2.0

GOLANG_VERSION = 1.14

.PHONY: clear
clear: ## Clear the working area and the project
	rm -rf bin/

.PHONY: run-%
run-%: build-%
	${BUILD_DIR}/$*

.PHONY: run
run: $(patsubst cmd/%,run-%,$(wildcard cmd/*)) ## Build and execute a binary

.PHONY: clean
clean: ## Clean builds
	rm -rf ${BUILD_DIR}/
	rm -rf cmd/*/pkged.go

.PHONY: goversion
goversion:
ifneq (${IGNORE_GOLANG_VERSION_REQ}, 1)
	@printf "${GOLANG_VERSION}\n$$(go version | awk '{sub(/^go/, "", $$3);print $$3}')" | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -g | head -1 | grep -q -E "^${GOLANG_VERSION}$$" || (printf "Required Go version is ${GOLANG_VERSION}\nInstalled: `go version`" && exit 1)
endif

.PHONY: build-%
build-%: goversion
ifeq (${VERBOSE}, 1)
	go env
endif

	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/$* ./cmd/$*

.PHONY: build
build: goversion ## Build all binaries
ifeq (${VERBOSE}, 1)
	go env
endif

	@mkdir -p ${BUILD_DIR}
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/ ./cmd/...

.PHONY: build-release
build-release: $(patsubst cmd/%,cmd/%/pkged.go,$(wildcard cmd/*)) ## Build all binaries without debug information
	@${MAKE} LDFLAGS="-w ${LDFLAGS}" GOARGS="${GOARGS} -trimpath" BUILD_DIR="${BUILD_DIR}/release" build

.PHONY: build-debug
build-debug: ## Build all binaries with remote debugging capabilities
	@${MAKE} GOARGS="${GOARGS} -gcflags \"all=-N -l\"" BUILD_DIR="${BUILD_DIR}/debug" build

bin/pkger:
	@mkdir -p bin
	go build -o bin/pkger github.com/markbates/pkger/cmd/pkger

cmd/%/pkged.go: bin/pkger ## Embed static files
	bin/pkger -o cmd/$*

.PHONY: check
check: test lint ## Run tests and linters

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum
bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_${OS}_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}

TEST_PKGS ?= ./...
.PHONY: test
test: TEST_FORMAT ?= short
test: SHELL = /bin/bash
test: export CGO_ENABLED=1
test: bin/gotestsum ## Run tests
	@mkdir -p ${BUILD_DIR}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/coverage.xml --format ${TEST_FORMAT} -- -race -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run

.PHONY: fix
fix: bin/golangci-lint ## Fix lint violations
	bin/golangci-lint run --fix

bin/mga: bin/mga-${MGA_VERSION}
	@ln -sf mga-${MGA_VERSION} bin/mga
bin/mga-${MGA_VERSION}:
	@mkdir -p bin
	curl -sfL https://git.io/mgatool | bash -s v${MGA_VERSION}
	@mv bin/mga $@

.PHONY: generate
generate: bin/mga ## Generate code
	go generate -x ./...
	mga generate kit endpoint ./...

.PHONY: openapi
openapi: ## Generate client and server stubs from the OpenAPI definition
	rm -rf internal/.generated/api/v1/rest
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v${OPENAPI_GENERATOR_VERSION} generate \
	--additional-properties packageName=api \
	--additional-properties withGoCodegenComment=true \
	-i /local/api/todo/v1/openapi.yaml \
	-g go-server \
	-o /local/internal/.generated/api/v1/rest
	rm -rf internal/.generated/api/v1/rest/{Dockerfile,go.*,README.md,main.go,go/api*.go,go/logger.go,go/routers.go}

	rm -rf api/todo/v1/client/rest
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v${OPENAPI_GENERATOR_VERSION} generate \
	--additional-properties packageName=todov1 \
	--additional-properties withGoCodegenComment=true \
	-i /local/api/todo/v1/openapi.yaml \
	-g go \
	-o /local/api/todo/v1/client/rest
	sed 's#jsonCheck = .*#jsonCheck = regexp.MustCompile(`(?i:(?:application|text)/(?:(?:vnd\\.[^;]+\\+)|(?:problem\\+))?json)`)#' api/todo/v1/client/rest/client.go > api/todo/v1/client/rest/client.go.new
	mv api/todo/v1/client/rest/client.go.new api/todo/v1/client/rest/client.go
	rm api/todo/v1/client/rest/{.travis.yml,git_push.sh,go.*}

bin/protoc: bin/protoc-${PROTOC_VERSION}
	@ln -sf protoc-${PROTOC_VERSION}/bin/protoc bin/protoc
bin/protoc-${PROTOC_VERSION}:
	@mkdir -p bin/protoc-${PROTOC_VERSION}
ifeq (${OS}, darwin)
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-osx-x86_64.zip > bin/protoc.zip
endif
ifeq (${OS}, linux)
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip > bin/protoc.zip
endif
	unzip bin/protoc.zip -d bin/protoc-${PROTOC_VERSION}
	rm bin/protoc.zip

bin/protoc-gen-go: go.mod
	@mkdir -p bin
	go build -o bin/protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go

bin/protoc-gen-go-grpc: gotools.mod
	@mkdir -p bin
	go build -modfile gotools.mod -o bin/protoc-gen-go-grpc google.golang.org/grpc/cmd/protoc-gen-go-grpc

bin/protoc-gen-kit: bin/protoc-gen-kit-${PROTOC_GEN_KIT_VERSION}
	@ln -sf protoc-gen-kit-${PROTOC_GEN_KIT_VERSION} bin/protoc-gen-kit
bin/protoc-gen-kit-${PROTOC_GEN_KIT_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/sagikazarmark/protoc-gen-kit/releases/download/v${PROTOC_GEN_KIT_VERSION}/protoc-gen-kit_${OS}_amd64.tar.gz | tar -zOxf - protoc-gen-kit > ./bin/protoc-gen-kit-${PROTOC_GEN_KIT_VERSION} && chmod +x ./bin/protoc-gen-kit-${PROTOC_GEN_KIT_VERSION}

bin/buf: bin/buf-${BUF_VERSION}
	@ln -sf buf-${BUF_VERSION} bin/buf
bin/buf-${BUF_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-${OS}-x86_64 -o ./bin/buf-${BUF_VERSION} && chmod +x ./bin/buf-${BUF_VERSION}

.PHONY: buf
buf: bin/buf ## Generate client and server stubs from the protobuf definition
	buf image build -o /dev/null
	buf check lint

.PHONY: proto
proto: bin/protoc bin/protoc-gen-go bin/protoc-gen-go-grpc bin/protoc-gen-kit buf ## Generate client and server stubs from the protobuf definition
	buf image build -o - | protoc --descriptor_set_in=/dev/stdin --go_out=paths=source_relative:api --go-grpc_out=paths=source_relative:api --kit_out=paths=source_relative:api $(shell buf image build -o - | buf ls-files --input - | grep -v google)

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)
