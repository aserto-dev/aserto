package dev

import (
	"path"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"
	"github.com/aserto-dev/aserto/pkg/filex"
	localpaths "github.com/aserto-dev/aserto/pkg/paths"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

const (
	local = "local"
)

type StartCmd struct {
	Name             string `arg:"" required:"" help:"policy name"`
	SrcPath          string `optional:"" type:"path" help:"path to source or bundle file"`
	Interactive      bool   `optional:""  help:"interactive execution mode instead of the default daemon mode"`
	ContainerName    string `optional:""  default:"authorizer-onebox" help:"container name"`
	ContainerVersion string `optional:""  default:"latest" help:"container version" `
}

func (cmd *StartCmd) Run(c *cc.CommonCtx) error {
	if running, err := dockerx.IsRunning(dockerx.AsertoOne); running || err != nil {
		if err != nil {
			return err
		}
		color.Yellow("!!! onebox is already running")
		return nil
	}

	color.Green(">>> starting onebox...")

	paths, err := localpaths.New()
	if err != nil {
		return err
	}

	if cmd.Name == local {
		if err := verifySrcPath(cmd.SrcPath); err != nil {
			return err
		}
	} else if !filex.FileExists(path.Join(paths.Config, cmd.Name+".yaml")) {
		return errors.Errorf("config for policy [%s] not found\nplease ensure the name is correct or\n run \"aserto developer configure <name>\" to create or update the policy configuration file", cmd.Name)
	}

	args := cmd.dockerArgs()

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

	return dockerx.DockerWith(cmd.env(paths), args...)
}

var (
	dockerCmd = []string{
		"run",
	}

	dockerArgs = []string{
		"--rm",
		"--name", dockerx.AsertoOne,
		"--platform", "linux/amd64",
		"-p", "8282:8282",
		"-p", "8383:8383",
		"-p", "8484:8484",
		"-v", "$ASERTO_CERTS_DIR:/certs:rw",
		"-v", "$ASERTO_CFG_DIR:/app/cfg:ro",
		"-v", "$ASERTO_EDS_DIR:/app/db:rw",
	}

	interactiveArgs = []string{
		"-ti",
	}

	daemonArgs = []string{
		"-d",
	}

	srcVolume = []string{
		"-v", "$ASERTO_SRC_DIR:/app/src:rw",
	}

	containerName = []string{
		"ghcr.io/aserto-dev/$CONTAINER_NAME:$CONTAINER_VERSION",
	}
)

func (cmd *StartCmd) dockerArgs() []string {
	args := append([]string{}, dockerCmd...)
	args = append(args, dockerArgs...)

	if cmd.Interactive {
		args = append(args, interactiveArgs...)
	} else {
		args = append(args, daemonArgs...)
	}

	if cmd.Name == local {
		args = append(args, srcVolume...)
	}

	return append(args, containerName...)
}

func (cmd *StartCmd) env(paths *localpaths.Paths) map[string]string {
	return map[string]string{
		"ASERTO_CERTS_DIR":  paths.Certs.Root,
		"ASERTO_CFG_DIR":    paths.Config,
		"ASERTO_EDS_DIR":    paths.EDS,
		"ASERTO_SRC_DIR":    cmd.SrcPath,
		"CONTAINER_NAME":    cmd.ContainerName,
		"CONTAINER_VERSION": cmd.ContainerVersion,
	}
}

func verifySrcPath(srcPath string) error {
	if srcPath == "" {
		return errors.Errorf("mode local requires source path argument to be set")
	}

	if !filex.DirExists(srcPath) {
		return errors.Errorf("source path directory %s does not exist", srcPath)
	}

	return nil
}
