filegroup(
    name = "go-deps",
    exported_deps = [
        "//third_party/go:github.com__golang__protobuf__proto",
        "//third_party/go:google.golang.org__protobuf__reflect__protoreflect",
        "//third_party/go:google.golang.org__protobuf__runtime__protoimpl",
        "//third_party/go:google.golang.org__protobuf__types__known__wrapperspb",
    ],
    needs_transitive_deps = True,
    output_is_complete = False,
    requires = ["go"],
    visibility = ["PUBLIC"],
)

filegroup(
    name = "go-grpc-deps",
    exported_deps = [
        "//third_party/go:google.golang.org__grpc",
        "//third_party/go:google.golang.org__grpc__codes",
        "//third_party/go:google.golang.org__grpc__status",
    ],
    needs_transitive_deps = True,
    output_is_complete = False,
    requires = ["go"],
    visibility = ["PUBLIC"],
)
