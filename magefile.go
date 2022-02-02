//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
)

func init() {
	// Set private repositories
	os.Setenv("GOPRIVATE", "github.com/aserto-dev")
}

// Lint runs linting for the entire project.
func Lint() error {
	return common.Lint()
}

// Test runs all tests and generates a code coverage report.
func Test() error {
	return common.Test()
}

// Build builds all binaries in ./cmd.
func Build() error {
	return common.BuildReleaser()
}

// Build and publish to GitHub.
func Publish() error {
	return common.Release("--rm-dist", "--snapshot", "--config", ".goreleaser-prod.yml")
}

// Release releases the project.
func Release() error {
	return common.Release("--rm-dist", "--config", ".goreleaser-prod.yml")
}

// BuildAll builds all binaries in ./cmd for
// all configured operating systems and architectures.
func BuildAll() error {
	return common.BuildAllReleaser()
}

// Deps installs all dependencies required to build the project.
func Deps() {
	deps.GetAllDeps()
}
