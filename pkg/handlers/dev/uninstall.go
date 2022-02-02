package dev

import (
	"fmt"
	"os"
	"path"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/aserto/pkg/handlers/dev/certs"
	localpaths "github.com/aserto-dev/aserto/pkg/paths"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type UninstallCmd struct{}

func (cmd UninstallCmd) Run(c *cc.CommonCtx) error {
	color.Green(">>> uninstalling onebox...")

	var err error

	//nolint:gocritic
	if err = (StopCmd{}).Run(c); err != nil {
		return err
	}

	paths, err := localpaths.New()
	if err != nil {
		return errors.Wrap(err, "can't find aserto directories")

	}

	cfgLocal := path.Join(paths.Config, "local.yaml")
	if filex.FileExists(cfgLocal) {
		fmt.Fprintf(c.OutWriter, "removing %s\n", cfgLocal)
		if err = os.Remove(cfgLocal); err != nil {
			return errors.Wrapf(err, "removing %s", cfgLocal)
		}
	}

	edsFile := path.Join(paths.EDS, "eds-acmecorp-v4.db")
	if filex.FileExists(edsFile) {
		fmt.Fprintf(c.OutWriter, "removing %s\n", edsFile)
		if err = os.Remove(edsFile); err != nil {
			return errors.Wrapf(err, "removing %s", edsFile)
		}
	}

	if err = certs.RemoveTrustedCert(paths.Certs.Gateway.CA); err != nil {
		return errors.Wrap(err, "failed to remove trusted ca cert")
	}

	if err = os.RemoveAll(paths.Certs.Root); err != nil {
		return errors.Wrap(err, "failed to delete onebox certificates")
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
