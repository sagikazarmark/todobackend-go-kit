# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

export PATH := $(abspath bin/protoc/bin/):$(abspath bin/):${PATH}

# Build variables
BUILD_DIR ?= build
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
LDFLAGS += -X main.version=${VERSION}
export CGO_ENABLED ?= 0

.PHONY: build
build: ## Build all binaries
	@mkdir -p ${BUILD_DIR}
	go build -trimpath -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/ .

.PHONY: run
run: build ## Build and run the application
	${BUILD_DIR}/todobackend-go-kit

.PHONY: check
check: test lint ## Run checks (tests and linters)

.PHONY: test
test: TEST_FORMAT ?= short
test: export CGO_ENABLED=1
test: ## Run tests
	@mkdir -p ${BUILD_DIR}
	gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/coverage.xml --jsonfile ${BUILD_DIR}/test.json --format ${TEST_FORMAT} -- -race -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run ${LINT_ARGS}

.PHONY: fmt
fmt: ## Format code
	golangci-lint run --fix

.PHONY: proto
proto:
	protoc -I api/ --go_out=paths=source_relative:api/ --go-grpc_out=paths=source_relative:api/ --go-kit_out=paths=source_relative:api/ api/todo/v1/*.proto

.PHONY: graphql
graphql: ## Generate GraphQL code
	go run github.com/99designs/gqlgen generate

.PHONY: openapi
openapi: ## Generate go server based on openapi description
	openapi-generator-cli generate \
	--additional-properties packageName=api \
	--additional-properties sourceFolder=api \
	--additional-properties withGoCodegenComment=true \
	-i api/todo/v1/openapi.yaml \
	-g go-server \
	-o internal/generated/api/todo/v1/rest

# Dependency versions
GOTESTSUM_VERSION ?= 1.8.0
GOLANGCI_VERSION ?= 1.48.0
PROTOC_VERSION ?= 3.19.4
PROTOC_GEN_GO_VERSION ?= 1.28.1
PROTOC_GEN_GO_GRPC_VERSION ?= 1.2.0
PROTOC_GEN_GO_KIT_VERSION ?= 0.1.1
GQLGEN_VERSION ?= 0.17.8

deps: bin/gotestsum bin/golangci-lint bin/protoc bin/protoc-gen-go bin/protoc-gen-go-grpc bin/protoc-gen-go-kit bin/gqlgen

bin/gotestsum:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_$(shell uname | tr A-Z a-z)_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum
	@chmod +x ./bin/gotestsum

bin/golangci-lint:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}

bin/protoc:
	@mkdir -p bin/protoc
ifeq ($(shell uname | tr A-Z a-z), darwin)
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-osx-x86_64.zip > bin/protoc.zip
endif
ifeq ($(shell uname | tr A-Z a-z), linux)
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip > bin/protoc.zip
endif
	unzip bin/protoc.zip -d bin/protoc
	rm bin/protoc.zip

bin/protoc-gen-go:
	@mkdir -p bin
	curl -L https://github.com/protocolbuffers/protobuf-go/releases/download/v${PROTOC_GEN_GO_VERSION}/protoc-gen-go.v${PROTOC_GEN_GO_VERSION}.$(shell uname | tr A-Z a-z).amd64.tar.gz | tar -zOxf - protoc-gen-go > ./bin/protoc-gen-go
	@chmod +x ./bin/protoc-gen-go

bin/protoc-gen-go-grpc:
	@mkdir -p bin
	curl -L https://github.com/grpc/grpc-go/releases/download/cmd/protoc-gen-go-grpc/v${PROTOC_GEN_GO_GRPC_VERSION}/protoc-gen-go-grpc.v${PROTOC_GEN_GO_GRPC_VERSION}.$(shell uname | tr A-Z a-z).amd64.tar.gz | tar -zOxf - ./protoc-gen-go-grpc > ./bin/protoc-gen-go-grpc
	@chmod +x ./bin/protoc-gen-go-grpc

bin/protoc-gen-go-kit:
	@mkdir -p bin
	curl -L https://github.com/sagikazarmark/protoc-gen-go-kit/releases/download/v${PROTOC_GEN_GO_KIT_VERSION}/protoc-gen-go-kit_$(shell uname | tr A-Z a-z)_amd64.tar.gz | tar -zOxf - protoc-gen-go-kit > ./bin/protoc-gen-go-kit
	@chmod +x ./bin/protoc-gen-go-kit

bin/gqlgen:
	@mkdir -p bin
	GOBIN=${PWD}/bin/ go install github.com/99designs/gqlgen@v${GQLGEN_VERSION}

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'
