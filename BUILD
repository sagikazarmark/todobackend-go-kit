github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "f8a12721c6f929db3e227e07c152d428ac47ab1b",
)

sh_cmd(
    name = "generate",
    deps = ["///pleasings2//tools/go:mga"],
    cmd = [
        "$(out_exe ///pleasings2//tools/go:mga) generate kit endpoint ./...",
    ],
)
