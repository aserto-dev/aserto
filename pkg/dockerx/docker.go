package dockerx

import (
	"fmt"

	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

const (
	AsertoOne = "topaz"
)

var (
	DockerRun = sh.RunCmd("docker")
	DockerOut = sh.OutCmd("docker")
)

func DockerWith(env map[string]string, args ...string) error {
	return sh.RunWithV(env, "docker", args...)
}

func DockerWithOut(env map[string]string, args ...string) (string, error) {
	return sh.OutputWith(env, "docker", args...)
}

func IsRunning(name string) (bool, error) {
	if name == "" {
		return false, errors.Errorf("instance name not specified")
	}
	str, err := DockerOut("ps", "-q", "-f", fmt.Sprintf("name=%s", name))
	return str != "", err
}
