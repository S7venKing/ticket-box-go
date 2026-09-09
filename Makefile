# Root Makefile. Recipe lines start with ">" (no tab required).
.RECIPEPREFIX = >

BUF_VERSION                ?= v1.72.0
PROTOC_GEN_GO_VERSION      ?= v1.36.12
PROTOC_GEN_GO_GRPC_VERSION ?= v1.6.2

MODULE_PREFIX := github.com/S7venKing/ticket-box-go

.PHONY: tools proto-lint proto-breaking proto-gen tidy sync build vet test

tools:
> go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
> go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
> go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

proto-lint:
> buf lint

proto-breaking:
> buf breaking --against '.git#branch=main,subdir=proto'

proto-gen:
> buf generate

tidy:
> cd gen && go mod tidy
> cd services/identity-service && go mod tidy

sync:
> go work sync

# In workspace mode "./..." does not work at the root (root is not a module);
# the module-path pattern covers every module listed in go.work.
build:
> go build $(MODULE_PREFIX)/...

vet:
> go vet $(MODULE_PREFIX)/...

test:
> go test $(MODULE_PREFIX)/...
