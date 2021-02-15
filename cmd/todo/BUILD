subinclude("///pleasings2//go")

timestamp = git_show("%ct")

date_fmt = "+%FT%T%z"

go_binary(
    name = "todo",
    srcs = glob(
        ["*.go"],
        exclude = [
            "*_test.go",
            "bindata.go",
        ],
    ) + [":bindata"],
    definitions = {
        "main.version": "${VERSION:-" + git_branch() + "}",
        "main.commitHash": git_commit()[0:8],
        "main.buildDate": f'$(date -u -d "@{timestamp}" "{date_fmt}" 2>/dev/null || date -u -r "{timestamp}" "{date_fmt}" 2>/dev/null || date -u "{date_fmt}")',
    },
    labels = ["binary"],
    #pass_env = ["VERSION"],
    #trimpath = True,
    visibility = ["PUBLIC"],
    deps = [
        "//api/todo/v1",
        "//internal/.generated/api/v1/graphql",
        "//pkg/todo",
        "//pkg/todo/tododriver",
        "//third_party/go:github.com__99designs__gqlgen__graphql__handler",
        "//third_party/go:github.com__go-kit__kit__transport__http",
        "//third_party/go:github.com__goph__idgen__ulidgen",
        "//third_party/go:github.com__gorilla__handlers",
        "//third_party/go:github.com__gorilla__mux",
        "//third_party/go:github.com__oklog__run",
        "//third_party/go:github.com__spf13__pflag",
        "//third_party/go:google.golang.org__grpc",
    ],
)

go_bindata(
    name = "bindata",
    srcs = ["//static"],
)
