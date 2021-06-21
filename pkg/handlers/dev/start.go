package dev

import (
	"os"
	"path"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"
	"github.com/aserto-dev/aserto/pkg/filex"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

const (
	local = "local"
	// modeRemote = "remote"
)

type StartCmd struct {
	Name             string `arg:"" required:"" help:"policy name"`
	SrcPath          string `optional:"" type:"path" help:"path to source or bundle file"`
	Interactive      bool   `optional:""  help:"interactive execution mode instead of the default daemon mode"`
	ContainerName    string `optional:""  default:"aserto-one" help:"container name"`
	ContainerVersion string `optional:""  default:"latest" help:"container version" `
}

// nolint:funlen // tbd
func (cmd StartCmd) Run(c *cc.CommonCtx) error {
	if running, err := dockerx.IsRunning(dockerx.AsertoOne); running || err != nil {
		if err != nil {
			return err
		}
		color.Yellow("!!! onebox is already running")
		return nil
	}

	color.Green(">>> starting onebox...")

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if cmd.Name == local && cmd.SrcPath == "" {
		return errors.Errorf("mode local requires source path argument to be set")
	}

	if cmd.Name == local && !filex.DirExists(cmd.SrcPath) {
		return errors.Errorf("source path directory %s does not exist", cmd.SrcPath)
	}

	if cmd.Name != local && !filex.FileExists(path.Join(home, ".config/aserto/aserto-one/cfg", cmd.Name+".yaml")) {
		return errors.Errorf("config for policy [%s] not found\nplease ensure the name is correct or\n run \"aserto developer configure <name>\" to create or update the policy configuration file", cmd.Name)
	}

	env := map[string]string{
		"ASERTO_CERTS_DIR":  path.Join(home, ".config/aserto/aserto-one/certs"),
		"ASERTO_CFG_DIR":    path.Join(home, ".config/aserto/aserto-one/cfg"),
		"ASERTO_EDS_DIR":    path.Join(home, ".cache/aserto/aserto-one/eds"),
		"ASERTO_SRC_DIR":    cmd.SrcPath,
		"CONTAINER_NAME":    cmd.ContainerName,
		"CONTAINER_VERSION": cmd.ContainerVersion,
	}

	dockerCmd := []string{
		"run",
	}

	dockerArgs := []string{
		"--rm",
		"--name", "$CONTAINER_NAME",
		"--platform", "linux/amd64",
		"-p", "8282:8282",
		"-p", "8383:8383",
		"-p", "8484:8484",
		"-v", "$ASERTO_CERTS_DIR:/root/.config/aserto/aserto-one/certs/:rw",
		"-v", "$ASERTO_CFG_DIR:/app/cfg:ro",
		"-v", "$ASERTO_EDS_DIR:/app/eds:rw",
	}

	interactiveArgs := []string{
		"-ti",
	}

	daemonArgs := []string{
		"-d",
	}

	srcVolume := []string{
		"-v", "$ASERTO_SRC_DIR:/app/src:rw",
	}

	containerName := []string{
		"ghcr.io/aserto-dev/$CONTAINER_NAME:$CONTAINER_VERSION",
	}

	args := []string{}

	args = append(args, dockerCmd...)
	args = append(args, dockerArgs...)

	if cmd.Interactive {
		args = append(args, interactiveArgs...)
	} else {
		args = append(args, daemonArgs...)
	}

	if cmd.Name == local {
		args = append(args, srcVolume...)
	}

	args = append(args, containerName...)

	switch cmd.Name {
	case local:
		cmdArgs := []string{
			"run",
			"--config-file", "/app/cfg/local.yaml",
			"--bundle", "/app/src",
			"--watch",
		}

		args = append(args, cmdArgs...)
	default:
		cmdArgs := []string{
			"run",
			"--config-file", "/app/cfg/" + cmd.Name + ".yaml",
		}

		args = append(args, cmdArgs...)
	}

	return dockerx.DockerWith(env, args...)
}
