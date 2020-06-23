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

subinclude('//build_defs/go')

go_test2(
    name = 'test',
    build_in_tree = True,
)

build_rule(
    name = "wollemi",
    binary = True,
    srcs = [remote_file(
        name = "wollemi",
        _tag = "download",
        url = f"https://github.com/tcncloud/wollemi/releases/download/v0.0.3/wollemi-v0.0.3-{CONFIG.HOSTOS}-{CONFIG.HOSTARCH}.tar.gz"
    )],
    cmd = " && ".join([
        "tar xf $SRCS",
    ]),
    outs = ["wollemi"],
    visibility = ["PUBLIC"],
)
