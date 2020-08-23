package(default_visibility = ["PUBLIC"])

github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "fa8f66a0e9f22323769661865e0a5ececa12e050",
)

sh_cmd(
    name = "generate",
    deps = ["///pleasings2//tools/go:mga"],
    cmd = [
        "$(out_location ///pleasings2//tools/go:mga) generate kit endpoint ./...",
    ],
)
