subinclude("///pleasings2//go")

go_build(
    name = "todo",
    definitions = {
        "main.version": "${VERSION:-" + git_branch() + "}",
        "main.commitHash": git_commit()[0:8],
        "main.buildDate": '$($(echo $TMP_DIR | sed "s/plz-out\/.*//")/$(out_location :builddate))',
    },
    pass_env = ["VERSION", "SOURCE_DATE_EPOCH"],
    deps = [":builddate"],
    stamp = True,
    trimpath = True,
)

sh_cmd(
    name = "builddate",
    cmd = '\n'.join([
        'DATE_FMT="+%FT%T%z"',
        '',
        'if [[ "\\\$SOURCE_DATE_EPOCH" != "" ]]; then',
        '    date -u -d "@\\\${SOURCE_DATE_EPOCH}" "\\\${DATE_FMT}" 2>/dev/null || date -u -r "\\\${SOURCE_DATE_EPOCH}" "\\\${DATE_FMT}" 2>/dev/null || date -u "\\\${DATE_FMT}"',
        'else',
        '    date "\\\${DATE_FMT}"',
        'fi',
    ]),
)
