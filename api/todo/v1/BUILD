subinclude("///pleasings2//openapi")

genrule(
    name = "openapi_client",
    srcs = [openapi_library(
        name = "openapi_client_internal",
        out = "api/todo/v1/client/rest.internal",
        additional_properties = {
            "packageName": "todov1",
            "withGoCodegenComment": "true",
        },
        generator = "go",
        spec = "openapi.yaml",
        strip = ["{.travis.yml,git_push.sh,go.*}"],
    )],
    outs = ["api/todo/v1/client/rest"],
    cmd = [
        "sed 's#jsonCheck = .*#jsonCheck = regexp.MustCompile(`(?i:(?:application|text)/(?:(?:vnd\\\\.[^;]+\\\\+)|(?:problem\\\\+))?json)`)#' api/todo/v1/api/todo/v1/client/rest.internal/client.go > api/todo/v1/api/todo/v1/client/rest.internal/client.go.new",
        "mv api/todo/v1/api/todo/v1/client/rest.internal/client.go.new api/todo/v1/api/todo/v1/client/rest.internal/client.go",
        "mv api/todo/v1/api/todo/v1/client/rest.internal/ $OUT",
    ],
)

openapi_library(
    name = "openapi_server",
    out = "internal/.generated/api/v1/rest",
    additional_properties = {
        "packageName": "api",
        "withGoCodegenComment": "true",
    },
    generator = "go-server",
    spec = "openapi.yaml",
    strip = ["{Dockerfile,go.*,README.md,main.go,go/api*.go,go/logger.go,go/routers.go}"],
)

openapi_validate(
    name = "openapi_validate",
    spec = "openapi.yaml",
)

go_library(
    name = "v1",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/go:github.com__golang__protobuf__proto",
        "//third_party/go:github.com__golang__protobuf__ptypes__wrappers",
        "//third_party/go:google.golang.org__grpc",
        "//third_party/go:google.golang.org__grpc__codes",
        "//third_party/go:google.golang.org__grpc__status",
        "//third_party/go:google.golang.org__protobuf__reflect__protoreflect",
        "//third_party/go:google.golang.org__protobuf__runtime__protoimpl",
    ],
)
