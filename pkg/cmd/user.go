package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
)

type UserCmd struct {
	Info user.InfoCmd `cmd:"" group:"user" help:"get user profile information"`
	Get  user.GetCmd  `cmd:"" group:"user" help:"get property"`
}

func (cmd *UserCmd) BeforeApply(c *cc.CommonCtx) error {
	return c.VerifyLoggedIn()
}

func (cmd *UserCmd) Run(c *cc.CommonCtx) error {
	return nil
}
