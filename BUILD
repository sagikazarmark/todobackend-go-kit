package(default_visibility=['PUBLIC'])

github_repo(
    name = "please-go",
    repo = "sagikazarmark/please-go",
    revision = "master",
)

sh_cmd(
    name = 'generate',
    srcs = ['///please-go//tools/mga'],
    cmd = [
        '$(out_location ///please-go//tools/mga) generate kit endpoint ./...'
    ]
)
