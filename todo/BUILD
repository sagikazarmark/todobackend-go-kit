subinclude("///pleasings2//go")

filegroup(
    name = "todo",
    srcs = glob(["*.go"], exclude = ["*_test.go"]),
)

gotestsum_test(
    name = "todo_test",
    deps = [":todo", "//todo/features"],
    race = True,
)
