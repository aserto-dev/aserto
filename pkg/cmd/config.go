package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/config"
)

type ConfigCmd struct {
	GetTenant config.GetTenantCmd `cmd:"" help:"get tenant list" group:"config"`
	SetTenant config.SetTenantCmd `cmd:"" help:"set default tenant" group:"config"`
	GetEnv    config.GetEnvCmd    `cmd:"" help:"get environment info" group:"config"`
}

func (cmd *ConfigCmd) BeforeApply(c *CLI) error {
	c.RequireLogin()
	return nil
}

func (cmd *ConfigCmd) Run(c *cc.CommonCtx) error {
	return nil
}
