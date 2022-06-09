package main

import (
	"strings"

	"dagger.io/dagger"

	"universe.dagger.io/alpine"
	"universe.dagger.io/docker"
	"universe.dagger.io/docker/cli"
	"universe.dagger.io/go"

	"github.com/sagikazarmark/todobackend-go-kit/ci/golangci"
	"github.com/sagikazarmark/todobackend-go-kit/ci/superlinter"
	"github.com/sagikazarmark/todobackend-go-kit/ci/editorconfig"
	dockerlint "github.com/sagikazarmark/todobackend-go-kit/ci/docker"
)

dagger.#Plan & {
	client: filesystem: ".": read: exclude: [
		".github",
		"bin",
		"build",
		"tmp",
	]
	client: filesystem: "./build": write: contents: actions.build.debug.output
	client: network: "unix:///var/run/docker.sock": connect: dagger.#Socket

	actions: {
		_source: client.filesystem["."].read.contents

		build: {
			debug: go.#Build & {
				source:  _source
				package: "."
				os:      *client.platform.os | "linux"
				arch:    client.platform.arch

				env: {
					CGO_ENABLED: "0"
				}
			}

			release: {
				"linux/amd64": _

				[platform=string]: go.#Build & {
					source:  _source
					package: "."
					os:      strings.Split(platform, "/")[0]
					arch:    strings.Split(platform, "/")[1]

					ldflags: "-s -w"

					env: {
						CGO_ENABLED: "0"
					}
				}
			}

			"docker": docker.#Build & {
				steps: [
					alpine.#Build,
					docker.#Copy & {
						contents: actions.build.release."linux/amd64".output
						dest:     "/usr/local/bin"
					},
					docker.#Set & {
						config: cmd: ["todobackend-go-kit", "--http-addr", ":8000", "--public-url", "${PUBLIC_URL}"]
					},
				]
			}
		}

		local: {
			"import-image": cli.#Load & {
				image: actions.build."docker".output
				host:  client.network."unix:///var/run/docker.sock".connect
				tag:   "sagikazarmark/todobackend-go-kit:latest"
			}
		}

		checks: {
			test: {
				// Go unit tests
				unit: go.#Test & {
					source:  _source
					package: "./..."

					command: flags: "-race": true
				}
			}

			lint: {
				go: golangci.#Lint & {
					source:  _source
					version: "1.46"
				}
				"superlinter": superlinter.#Lint & {
					source: _source
				}
				"editorconfig": editorconfig.#Lint & {
					source: _source
				}
				"docker": dockerlint.#Lint & {
					source: _source
				}
			}
		}
	}
}
