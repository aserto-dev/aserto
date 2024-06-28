package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
)

type UserCmd struct {
	Get UserGetCmd `cmd:"" help:"get info | token |"`
}

func (cmd *UserCmd) BeforeApply(so ServiceOptions) error {
	so.RequireToken()
	return nil
}

func (cmd *UserCmd) Run(c *cc.CommonCtx) error {
	return nil
}

type UserGetCmd struct {
	Profile           user.InfoCmd `cmd:"" help:"user profile"`
	AccessToken       user.GetCmd  `cmd:"" help:"access token"`
	TenantID          user.GetCmd  `cmd:"" help:"tenant ID"`
	AuthorizerKey     user.GetCmd  `cmd:"" help:"authorizer key"`
	DirectoryReadKey  user.GetCmd  `cmd:"" help:"directory read key"`
	DirectoryWriteKey user.GetCmd  `cmd:"" help:"directory write key"`
	DiscoveryKey      user.GetCmd  `cmd:"" help:"discovery key"`
	RegistryReadKey   user.GetCmd  `cmd:"" help:"registry read key"`
	RegistryWriteKey  user.GetCmd  `cmd:"" help:"registry write key"`
	DecisionLogsKey   user.GetCmd  `cmd:"" help:"decision logs key"`
	Token             user.GetCmd  `cmd:"" help:"token" hidden:""`
}
