package dev

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"
	"github.com/aserto-dev/aserto/pkg/handlers/dev/certs"
	localpaths "github.com/aserto-dev/aserto/pkg/paths"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type InstallCmd struct {
	TrustCert bool `optional:"" default:"false" help:"add topaz certificate to the system's trusted CAs"`

	ContainerName    string `optional:""  default:"topaz" help:"container name"`
	ContainerVersion string `optional:""  default:"latest" help:"container version" `
}

func (cmd InstallCmd) Run(c *cc.CommonCtx) error {
	if running, err := dockerx.IsRunning(dockerx.AsertoOne); running || err != nil {
		if err != nil {
			return err
		}
		color.Yellow("!!! sidecar is already running")
		return nil
	}

	color.Green(">>> installing topaz...")

	paths, err := localpaths.Create()
	if err != nil {
		return errors.Wrap(err, "failed to create configuration directory")
	}

	// Create sidecar certs if none exist.
	if err := certs.GenerateCerts(c.UI.Output(), c.UI.Err(), paths.Certs.GRPC, paths.Certs.Gateway); err != nil {
		return errors.Wrap(err, "failed to create dev certificates")
	}

	if cmd.TrustCert {
		fmt.Fprintln(c.UI.Output(), "Adding developer certificate to system store. You may need to provide credentials.")
		if err := certs.AddTrustedCert(paths.Certs.Gateway.CA); err != nil {
			return errors.Wrap(err, "add gateway-ca cert to trusted certs")
		}
	}

	return dockerx.DockerWith(map[string]string{
		"CONTAINER_NAME":    cmd.ContainerName,
		"CONTAINER_VERSION": cmd.ContainerVersion,
	},
		"pull", "ghcr.io/aserto-dev/$CONTAINER_NAME:$CONTAINER_VERSION",
	)
}
