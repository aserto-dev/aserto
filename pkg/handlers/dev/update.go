package dev

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"

	"github.com/fatih/color"
)

type UpdateCmd struct {
	ContainerName    string `optional:""  default:"aserto-one" help:"container name"`
	ContainerVersion string `optional:""  default:"latest" help:"container version" `
}

func (cmd UpdateCmd) Run(c *cc.CommonCtx) error {
	color.Green(">>> updating onebox...")

	return dockerx.DockerWith(map[string]string{
		"CONTAINER_NAME":    cmd.ContainerName,
		"CONTAINER_VERSION": cmd.ContainerVersion,
	},
		"pull", "ghcr.io/aserto-dev/$CONTAINER_NAME:$CONTAINER_VERSION",
	)
}
