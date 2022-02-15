//go:build mage
// +build mage

package main

import (
	"fmt"
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
	return common.BuildReleaser("--config", ".goreleaser-prod.yml")
}

// Build and publish to GitHub.
func Publish() error {
	if os.Getenv("GITHUB_TOKEN") == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is undefined")
	}
	if os.Getenv("ASERTO_TAP") == "" {
		return fmt.Errorf("ASERTO_TAP environment variable is undefined")
	}

	return common.Release("--rm-dist", "--config", ".goreleaser-publish.yml")
}

// Release releases the project.
func Release() error {
	return common.Release("--skip-publish", "--rm-dist", "--snapshot", "--config", ".goreleaser-prod.yml")
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
