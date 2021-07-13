package cmd

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/version"
	"github.com/aserto-dev/aserto/pkg/x"
)

type CLI struct {
	Authorizer         AuthorizerCmd  `cmd:"" aliases:"a" help:"authorizer commands"`
	Tenant             TenantCmd      `cmd:"" aliases:"t" help:"tenant commands"`
	Directory          DirectoryCmd   `cmd:"" aliases:"d" help:"directory commands"`
	Developer          DeveloperCmd   `cmd:"" aliases:"x" help:"developer commands"`
	User               UserCmd        `cmd:"" aliases:"u" help:"user commands"`
	Login              user.LoginCmd  `cmd:"" help:"login"`
	Logout             user.LogoutCmd `cmd:"" help:"logout"`
	Config             ConfigCmd      `cmd:"" aliases:"c" help:"configuration commands"`
	Version            VersionCmd     `cmd:"" help:"version information"`
	Verbose            bool           `name:"verbose" help:"verbose output"`
	AuthorizerOverride string         `name:"authorizer" env:"ASERTO_AUTHORIZER" help:"authorizer override"`
	TenantOverride     string         `name:"tenant" env:"ASERTO_TENANT_ID" help:"tenant id override"`
	EnvOverride        string         `name:"env" default:"${defaultEnv}" env:"ASERTO_ENV" hidden:"" help:"environment override"`
	Debug              bool           `name:"debug" env:"ASERTO_DEBUG" help:"enable debug logging"`
	requireLogin       bool
}

func (cmd *CLI) RequireLogin() {
	cmd.requireLogin = true
}

func (cmd *CLI) IsLoginRequired() bool {
	return cmd.requireLogin
}

func (cmd *CLI) Run(c *cc.CommonCtx) error {
	return nil
}

type VersionCmd struct {
}

func (cmd *VersionCmd) Run(c *cc.CommonCtx) error {
	fmt.Fprintf(c.OutWriter, "%s - %s (%s)\n",
		x.AppName,
		version.GetInfo().String(),
		x.AppVersionTag,
	)
	return nil
}
