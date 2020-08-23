package(default_visibility = ["PUBLIC"])

github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "2248c153a30eb86e76ac1f55e7c4562425f6f145",
)

sh_cmd(
    name = "generate",
    deps = ["///pleasings2//tools/go:mga"],
    cmd = [
        "$(out_location ///pleasings2//tools/go:mga) generate kit endpoint ./...",
    ],
)
