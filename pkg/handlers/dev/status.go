package dev

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dockerx"

	"github.com/fatih/color"
)

type StatusCmd struct{}

func (cmd StatusCmd) Run(c *cc.CommonCtx) error {
	running, err := dockerx.IsRunning(dockerx.AsertoOne)
	if err != nil {
		return err
	}
	if running {
		color.Green(">>> topaz is running")
	} else {
		color.Yellow(">>> topaz is not running")
	}
	return nil
}
