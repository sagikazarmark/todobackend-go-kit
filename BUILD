github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "68f9d3727fcf5ee204288312682ceb51eff1eb83",
)

sh_cmd(
    name = "generate",
    deps = ["///pleasings2//tools/go:mga"],
    cmd = [
        "$(out_exe ///pleasings2//tools/go:mga) generate kit endpoint ./...",
    ],
)
