package superlinter

import (
	"dagger.io/dagger"
	"dagger.io/dagger/core"

	"universe.dagger.io/docker"
)

// Lint using golangci-lint
#Lint: {
	// Source code
	source: dagger.#FS

	// super-linter version
	version: *"4.9.4" | string

	// Environment variables
	env: [string]: string

	_image: docker.#Pull & {
		source: "github/super-linter:v\(version)"
	}

	docker.#Run & {
		input: *_image.output | docker.#Image
		mounts: {
			"source": {
				contents: source
				dest:     "/tmp/lint"
			}
			"golangci cache": {
				contents: core.#CacheDir & {
					id: "superlinter_golangci"
				}
				dest: "/root/.cache/golangci-lint"
			}
		}
		"env": {
			GOLANGCI_LINT_CACHE: "/root/.cache/golangci-lint"
			RUN_LOCAL:           "true"
			LOG_FILE:            "/dev/stdout"
		} & env
	}
}
