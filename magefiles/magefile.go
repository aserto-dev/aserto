//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
)

func init() {
	os.Setenv("GO_VERSION", "1.20")
	os.Setenv("DOCKER_BUILDKIT", "1")
}

// Generate generates all code.
func Generate() error {
	return common.Generate()
}

// Lint runs linting for the entire project.
func Lint() error {
	return common.Lint()
}

// Test runs all tests and generates a code coverage report.
func Test() error {
	return common.Test()
}

// Build all binaries in ./cmd.
func Build() error {
	return common.BuildReleaser()
}

// Release the project.
// func Release() error {
// 	if os.Getenv("GITHUB_TOKEN") == "" {
// 		return fmt.Errorf("GITHUB_TOKEN environment variable is undefined")
// 	}

// 	if os.Getenv("HOMEBREW_TAP") == "" {
// 		return fmt.Errorf("HOMEBREW_TAP environment variable is undefined")
// 	}

// 	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
// 		return fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is undefined")
// 	}

// 	if err := writeVersion(); err != nil {
// 		return err
// 	}

// 	return common.Release("--rm-dist")
// }

// BuildAll builds all binaries in ./cmd for
// all configured operating systems and architectures.
func BuildAll() error {
	return common.BuildAllReleaser("--rm-dist", "--snapshot")
}

// Deps installs all dependencies required to build the project.
func Deps() {
	deps.GetAllDeps()
}

// func writeVersion() error {
// 	version, err := exec.Command("git", "describe", "--tags").Output()
// 	if err != nil {
// 		return errors.Wrap(err, "failed to get current git tag")
// 	}

// 	file, err := os.Create("VERSION.txt")
// 	if err != nil {
// 		return errors.Wrap(err, "failed to create version file")
// 	}

// 	defer file.Close()

// 	if _, err := file.Write(version); err != nil {
// 		return errors.Wrap(err, "failed to write to version file")
// 	}

// 	return file.Sync()
// }
