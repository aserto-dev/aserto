package dev

import (
	"fmt"
	"os"
	"path"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"
	"github.com/aserto-dev/aserto/pkg/filex"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type UninstallCmd struct{}

func (cmd UninstallCmd) Run(c *cc.CommonCtx) error {
	color.Green(">>> uninstalling onebox...")

	var err error

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	//nolint:gocritic
	if err = (StopCmd{}).Run(c); err != nil {
		return err
	}

	cfgLocal := path.Join(home, ".config/aserto/aserto-one/cfg/local.yaml")
	if filex.FileExists(cfgLocal) {
		fmt.Fprintf(c.OutWriter, "removing %s\n", cfgLocal)
		if err = os.Remove(cfgLocal); err != nil {
			return errors.Wrapf(err, "removing %s", cfgLocal)
		}
	}

	edsFile := path.Join(home, ".cache/aserto/aserto-one/eds/eds-acmecorp-v4.db")
	if filex.FileExists(edsFile) {
		fmt.Fprintf(c.OutWriter, "removing %s\n", edsFile)
		if err = os.Remove(edsFile); err != nil {
			return errors.Wrapf(err, "removing %s", edsFile)
		}
	}

	str, err := dockerx.DockerWithOut(map[string]string{
		"NAME": "aserto-one",
	},
		"images",
		"ghcr.io/aserto-dev/$NAME",
		"--filter", "label=org.opencontainers.image.source=https://github.com/aserto-dev/aserto-one",
		"-q",
	)
	if err != nil {
		return err
	}

	if str != "" {
		fmt.Fprintf(c.OutWriter, "removing %s\n", "aserto-dev/aserto-one")
		err = dockerx.DockerRun("rmi", str)
	}

	return err
}
