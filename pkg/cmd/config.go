package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/config"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
)

type ConfigCmd struct {
	UserInfo    user.InfoCmd        `cmd:"" help:"get user profile information" group:"config"`
	GetProperty user.GetCmd         `cmd:"" help:"get property" group:"config"`
	GetTenant   config.GetTenantCmd `cmd:"" help:"get tenant list" group:"config"`
	SetTenant   config.SetTenantCmd `cmd:"" help:"set default tenant" group:"config"`
	GetEnv      config.GetEnvCmd    `cmd:"" help:"get environment info" group:"config"`
}

func (cmd *ConfigCmd) Run(c *cc.CommonCtx) error {
	return nil
}
