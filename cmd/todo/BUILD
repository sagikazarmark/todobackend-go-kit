subinclude('///please-go//rules')

go_build(
    name = 'todo',
    build_in_tree = True,
    stamp = True,
    trimpath = True,
    pass_env = ['VERSION'],
    definitions = {
        'main.version': '${VERSION:-' + git_branch() + '}',
        'main.commitHash': git_commit()[0:8],
        'main.buildDate': '$SCM_COMMIT_DATE',
    }
)
