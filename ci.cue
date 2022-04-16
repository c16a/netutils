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

		build: go.#Build & {
			source:  _source
			package: "github.com/c16a/netutils"
			os:      client.platform.os
			arch:    client.platform.arch

			ldflags: "-s -w"

			env: {
				CGO_ENABLED: "0"
				// Makes sure the linter and unit tests complete before starting the build
				// "__depends_lint":  "\(goLint.exit)"
				// "__depends_tests": "\(goTest.exit)"
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