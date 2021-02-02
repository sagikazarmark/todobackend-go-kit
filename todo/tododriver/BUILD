subinclude("///pleasings2//go")

filegroup(
    name = "tododriver",
    srcs = glob(["*.go"], exclude = ["*_test.go"]),
)

gotestsum_test(
    name = "tododriver_test",
    deps = [":tododriver"],
    race = True,
)
