github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "da9b05cf51fbaa900ddea8bac0739ee74df364bd",
)

sh_cmd(
    name = "generate",
    deps = ["///pleasings2//tools/go:mga"],
    cmd = [
        "$(out_exe ///pleasings2//tools/go:mga) generate kit endpoint ./...",
    ],
)

sh_cmd(
    name = "proto",
    deps = [
        "///pleasings2//tools/proto:buf",
        "///pleasings2//tools/proto:protoc",
        "///pleasings2//tools/proto:protoc-gen-go",
        "///pleasings2//tools/proto:protoc-gen-go-grpc",
        "///pleasings2//tools/proto:protoc-gen-kit",
    ],
    cmd = [
        "$(out_exe ///pleasings2//tools/proto:buf) image build -o - | $(out_exe ///pleasings2//tools/proto:protoc) --descriptor_set_in=/dev/stdin --plugin=$(out_exe ///pleasings2//tools/proto:protoc-gen-go) --plugin=$(out_exe ///pleasings2//tools/proto:protoc-gen-go-grpc) --plugin=$(out_exe ///pleasings2//tools/proto:protoc-gen-kit) --go_out=paths=source_relative:api --go-grpc_out=paths=source_relative:api --kit_out=paths=source_relative:api \\\$($(out_exe ///pleasings2//tools/proto:buf) image build -o - | buf ls-files --input - | grep -v google)",
    ],
)
