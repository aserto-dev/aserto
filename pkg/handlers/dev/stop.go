package dev

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"

	"github.com/fatih/color"
)

type StopCmd struct{}

func (cmd StopCmd) Run(c *cc.CommonCtx) error {
	running, err := dockerx.IsRunning(dockerx.AsertoOne)
	if err != nil {
		return err
	}

	if running {
		color.Green(">>> stopping onebox...")
		return dockerx.DockerRun("stop", "aserto-one")
	}

	return nil
}
