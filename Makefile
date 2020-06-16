# Hooks
CLEAN_TARGETS = clean-pkged
BUILD_RELEASE_DEP_TARGETS = bin/pkger
PRE_BUILD_RELEASE_TARGETS = $(patsubst cmd/%,cmd/%/pkged.go,$(wildcard cmd/*))
BUILD_DEBUG_DEP_TARGETS = bin/pkger
PRE_BUILD_DEBUG_TARGETS = $(patsubst cmd/%,cmd/%/pkged.go,$(wildcard cmd/*))

include main.mk

# Dependency versions
MGA_VERSION = 0.2.0
OPENAPI_GENERATOR_VERSION = 4.3.1
PROTOC_VERSION = 3.12.2
BUF_VERSION = 0.15.0
PROTOC_GEN_KIT_VERSION = 0.2.0

bin/pkger: go.mod
	@mkdir -p bin
	go build -o bin/pkger github.com/markbates/pkger/cmd/pkger

cmd/%/pkged.go: bin/pkger
	bin/pkger -o cmd/$*

.PHONY: clean-pkged
clean-pkged:
	rm -rf cmd/*/pkged.go

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

bin/gqlgen: go.mod
	@mkdir -p bin
	go build -o bin/gqlgen github.com/99designs/gqlgen

.PHONY: graphql
graphql: bin/gqlgen ## Generate GraphQL code
	bin/gqlgen
