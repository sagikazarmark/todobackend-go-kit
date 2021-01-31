github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "f6b0609fe9b4b406906e76a2e5c140017b2e3628",
)

sh_cmd(
    name = "generate",
    cmd = [
        "$(out_exe ///pleasings2//tools/go:mga) generate kit endpoint ./...",
    ],
    deps = ["///pleasings2//tools/go:mga"],
)

sh_cmd(
    name = "proto",
    cmd = [
        "$(out_exe ///pleasings2//tools/proto:buf) image build -o - | $(out_exe ///pleasings2//tools/proto:protoc) --descriptor_set_in=/dev/stdin --plugin=$(out_exe ///pleasings2//tools/proto:protoc-gen-go) --plugin=$(out_exe ///pleasings2//tools/proto:protoc-gen-go-grpc) --plugin=$(out_exe ///pleasings2//tools/proto:protoc-gen-kit) --go_out=paths=source_relative:api --go-grpc_out=paths=source_relative:api --kit_out=paths=source_relative:api \\\\$($(out_exe ///pleasings2//tools/proto:buf) image build -o - | buf ls-files --input - | grep -v google)",
    ],
    deps = [
        "///pleasings2//tools/proto:buf",
        "///pleasings2//tools/proto:protoc",
        "///pleasings2//tools/proto:protoc-gen-go",
        "///pleasings2//tools/proto:protoc-gen-go-grpc",
        "///pleasings2//tools/proto:protoc-gen-kit",
    ],
)
