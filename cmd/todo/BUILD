subinclude("///pleasings2//go")

timestamp = git_show("%ct")
date_fmt = "+%FT%T%z"

go_build(
    name = "todo",
    definitions = {
        "main.version": "${VERSION:-" + git_branch() + "}",
        "main.commitHash": git_commit()[0:8],
        "main.buildDate": f'$(date -u -d "@{timestamp}" "{date_fmt}" 2>/dev/null || date -u -r "{timestamp}" "{date_fmt}" 2>/dev/null || date -u "{date_fmt}")',
    },
    pass_env = ["VERSION"],
    stamp = True,
    trimpath = True,
)
