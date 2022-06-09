package editorconfig

import (
	"dagger.io/dagger"

	"universe.dagger.io/docker"
)

// Lint editorconfig files
#Lint: {
	// Source code
	source: dagger.#FS

	// super-linter version
	version: *"2.4.0" | string

	_image: docker.#Pull & {
		source: "docker.io/mstruebing/editorconfig-checker:\(version)"
	}

	_sourcePath: "/src"

	docker.#Run & {
		input:   *_image.output | docker.#Image
		workdir: _sourcePath
		mounts: {
			"source": {
				contents: source
				dest:     _sourcePath
			}
		}
	}
}
