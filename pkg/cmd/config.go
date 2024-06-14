package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/config"
)

type ConfigCmd struct {
	GetEnv config.GetEnvCmd     `cmd:"" help:"get environment info" group:"config"`
	Use    config.UseConfigCmd  `cmd:"" help:"use a topaz configuration" group:"config"`
	List   config.ListConfigCmd `cmd:"" help:"list configurations" group:"config"`
}

func (cmd *ConfigCmd) Run(c *cc.CommonCtx) error {
	return nil
}
