github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "ab4e4442d8106573e392d25545ce842d6d5e514c",
)

sh_cmd(
    name = "generate",
    deps = ["///pleasings2//tools/go:mga"],
    cmd = [
        "$(out_exe ///pleasings2//tools/go:mga) generate kit endpoint ./...",
    ],
)
