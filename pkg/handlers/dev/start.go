package dev

import (
	"fmt"
	"os"
	"path"

	"github.com/aserto-dev/aserto/pkg/cc"
	decisionlogger "github.com/aserto-dev/aserto/pkg/decision_logger"
	"github.com/aserto-dev/aserto/pkg/dockerx"
	"github.com/aserto-dev/aserto/pkg/filex"
	localpaths "github.com/aserto-dev/aserto/pkg/paths"
	"github.com/spf13/viper"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

const (
	local = "local"
)

type StartCmd struct {
	Name             string `arg:"" required:"" help:"policy name"`
	SrcPath          string `optional:"" type:"path" help:"path to source or bundle file"`
	Interactive      bool   `optional:"" help:"interactive execution mode instead of the default daemon mode"`
	ContainerName    string `optional:"" default:"topaz" help:"container name"`
	ContainerVersion string `optional:"" default:"latest" help:"container version" `
	Hostname         string `optional:"" help:"hostname for docker to set"`
	DataPath         string `optional:"" type:"path" help:"path for non-ephemeral data storage"`
}

func (cmd *StartCmd) Run(c *cc.CommonCtx) error {
	if running, err := dockerx.IsRunning(dockerx.AsertoOne); running || err != nil {
		if err != nil {
			return err
		}
		color.Yellow("!!! topaz is already running")
		return nil
	}

	color.Green(">>> starting topaz ...")

	paths, err := localpaths.NewWithDataRoot(cmd.DataPath)
	if err != nil {
		return err
	}

	err = cmd.validateConfig(c, paths)
	if err != nil {
		return err
	}

	args := cmd.dockerArgs()

	switch cmd.Name {
	case local:
		cmdArgs := []string{
			"run",
			"--config-file", "/cfg/local.yaml",
			"--bundle", "/app/src",
			"--watch",
		}

		args = append(args, cmdArgs...)
	default:
		cmdArgs := []string{
			"run",
			"--config-file", "/cfg/" + cmd.Name + ".yaml",
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
		"-p", "9292:9292",
		"-v", "$ASERTO_CERTS_DIR:/certs:rw",
		"-v", "$ASERTO_CFG_DIR:/cfg:ro",
		"-v", "$ASERTO_EDS_DIR:/db:rw",
		"-v", "$ASERTO_DECISION_LOGS_DIR:/app/decision_logs:rw",
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

	hostname = []string{
		"--hostname", "$CONTAINER_HOSTNAME",
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

	if cmd.Hostname != "" {
		args = append(args, hostname...)
	}

	return append(args, containerName...)
}

func (cmd *StartCmd) env(paths *localpaths.Paths) map[string]string {
	return map[string]string{
		"ASERTO_CERTS_DIR":         paths.Certs.Root,
		"ASERTO_CFG_DIR":           paths.Config,
		"ASERTO_EDS_DIR":           paths.EDS,
		"ASERTO_SRC_DIR":           cmd.SrcPath,
		"CONTAINER_NAME":           cmd.ContainerName,
		"CONTAINER_VERSION":        cmd.ContainerVersion,
		"CONTAINER_HOSTNAME":       cmd.Hostname,
		"ASERTO_DECISION_LOGS_DIR": path.Join(paths.Data, decisionlogger.Dir),
	}
}

func setupLocalRun(c *cc.CommonCtx, paths *localpaths.Paths, srcPath string) error {
	if srcPath == "" {
		return errors.Errorf("mode local requires source path argument to be set")
	}

	if !filex.DirExists(srcPath) {
		return errors.Errorf("source path directory %s does not exist", srcPath)
	}

	cfgLocal := paths.LocalConfig()
	if !filex.FileExists(cfgLocal) {
		fmt.Fprintf(c.UI.Output(), "creating %s\n", cfgLocal)
		f, err := os.OpenFile(cfgLocal, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return errors.Wrapf(err, "creating %s", cfgLocal)
		}

		if err := WriteConfig(f, configTemplateLocal, &templateParams{TenantID: c.TenantID()}); err != nil {
			return errors.Wrapf(err, "writing %s", cfgLocal)
		}

	}

	return nil
}

func (cmd *StartCmd) validateConfig(c *cc.CommonCtx, paths *localpaths.Paths) error {
	if cmd.Name == local {
		if err := setupLocalRun(c, paths, cmd.SrcPath); err != nil {
			return err
		}
	} else if !filex.FileExists(path.Join(paths.Config, cmd.Name+".yaml")) {
		return errors.Errorf("config for policy [%s] not found\nplease ensure the name is correct or\n run \"aserto developer configure <name>\" to create or update the policy configuration file", cmd.Name)
	}

	cfgPath := path.Join(paths.Config, cmd.Name+".yaml")
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(cfgPath)

	err := v.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "error reading config file '%s'", cfgPath)
	}

	cfg := v.GetStringMap("decision_logger")
	if cfg != nil && cmd.DataPath == "" {
		return errors.Errorf("policy '%s' has decision logging configured, please specify a destination for logs using --data-path", cmd.Name)
	}

	return nil
}
