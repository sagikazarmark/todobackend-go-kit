package(default_visibility = ["PUBLIC"])

github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "f845565009a33200a3b92404867bcba7148f758f",
)

sh_cmd(
    name = "generate",
    deps = ["///pleasings2//tools/go:mga"],
    cmd = [
        "$(out_location ///pleasings2//tools/go:mga) generate kit endpoint ./...",
    ],
)
