package docker

import (
	"dagger.io/dagger"

	"universe.dagger.io/docker"
)

// Lint editorconfig files
#Lint: {
	// Source code
	source: dagger.#FS

	dockerfile: *"Dockerfile" | string

	// super-linter version
	version: *"2.10.0" | string

	_image: docker.#Pull & {
		source: "docker.io/hadolint/hadolint:v\(version)"
	}

	_sourcePath: "/src"

	docker.#Run & {
		input:   *_image.output | docker.#Image
		workdir: _sourcePath
		command: {
			name: "hadolint"
			args: [dockerfile]
		}
		mounts: {
			"source": {
				contents: source
				dest:     _sourcePath
			}
		}
	}
}
