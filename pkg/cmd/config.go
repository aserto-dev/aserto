package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/config"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
)

type ConfigCmd struct {
	UserInfo         user.InfoCmd               `cmd:"" help:"get user profile information" group:"config"`
	GetProperty      user.GetCmd                `cmd:"" help:"get property" group:"config"`
	GetContexts      config.GetContextsCmd      `cmd:"" help:"get defined contexts" group:"config"`
	GetActiveContext config.GetActiveContextCmd `cmd:"" help:"get active context config" group:"config"`
	DeleteContext    config.DeleteContextCmd    `cmd:"" help:"delete a context config" group:"config"`
	SetContext       config.SetContextCmd       `cmd:"" help:"creates a context" group:"config"`
	UseContext       config.UseContextCmd       `cmd:"" help:"use a specific context" group:"config"`
}

func (cmd *ConfigCmd) Run(c *cc.CommonCtx) error {
	return nil
}
