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
