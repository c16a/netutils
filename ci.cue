package main

import (
	"dagger.io/dagger"
	"universe.dagger.io/go"
)

dagger.#Plan & {

	client: filesystem: ".": read: exclude: [
		"bin",
		"**/node_modules",
		"cmd/dagger/dagger",
		"cmd/dagger/dagger-debug",
	]
	client: filesystem: "./bin": write: contents: actions.build.output

	actions: {
		_source: client.filesystem["."].read.contents

		buildLinuxAmd64: go.#Build & {
			source:  _source
			package: "github.com/c16a/netutils"
			os:      "linux"
			arch:    "amd64"

			ldflags: "-s -w"

			env: {
				CGO_ENABLED: "0"
			}
		}

        buildLinuxArm64: go.#Build & {
			source:  _source
			package: "github.com/c16a/netutils"
			os:      "linux"
			arch:    "arm64"

			ldflags: "-s -w"

			env: {
				CGO_ENABLED: "0"
			}
		}

        buildDarwinAmd64: go.#Build & {
			source:  _source
			package: "github.com/c16a/netutils"
			os:      "darwin"
			arch:    "amd64"

			ldflags: "-s -w"

			env: {
				CGO_ENABLED: "0"
			}
		}

        buildDarwinArm64: go.#Build & {
			source:  _source
			package: "github.com/c16a/netutils"
			os:      "darwin"
			arch:    "arm64"

			ldflags: "-s -w"

			env: {
				CGO_ENABLED: "0"
			}
		}

        buildWindowsAmd64: go.#Build & {
			source:  _source
			package: "github.com/c16a/netutils"
			os:      "windows"
			arch:    "amd64"

			ldflags: "-s -w"

			env: {
				CGO_ENABLED: "0"
			}
		}

		// Go unit tests
		test: go.#Test & {
			// container: image: _goImage.output
			source:  _source
			package: "./..."

			// FIXME: doesn't work with CGO_ENABLED=0
			// command: flags: "-race": true

			env: {
				// FIXME: removing this complains about lack of gcc
				CGO_ENABLED: "0"
			}
		}
	}
}