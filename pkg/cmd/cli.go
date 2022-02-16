package cmd

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/version"
	"github.com/aserto-dev/aserto/pkg/x"
)

type CLI struct {
	Authorizer   AuthorizerCmd   `cmd:"" aliases:"a" help:"authorizer commands"`
	Tenant       TenantCmd       `cmd:"" aliases:"t" help:"tenant commands"`
	Directory    DirectoryCmd    `cmd:"" aliases:"d" help:"directory commands"`
	DecisionLogs DecisionLogsCmd `cmd:"" aliases:"l" help:"decision logs commands"`
	Developer    DeveloperCmd    `cmd:"" aliases:"x" help:"developer commands"`
	User         UserCmd         `cmd:"" aliases:"u" help:"user commands"`
	Login        user.LoginCmd   `cmd:"" help:"login"`
	Logout       user.LogoutCmd  `cmd:"" help:"logout"`
	Config       ConfigCmd       `cmd:"" aliases:"c" help:"configuration commands"`
	Version      VersionCmd      `cmd:"" help:"version information"`

	Debug          bool   `name:"debug" env:"ASERTO_DEBUG" help:"enable debug logging"`
	EnvOverride    string `name:"env" default:"${defaultEnv}" env:"ASERTO_ENV" hidden:"" help:"environment override"`
	TenantOverride string `name:"tenant" env:"ASERTO_TENANT_ID" help:"tenant id override"`
}

func (cli *CLI) TenantID(token cc.CachedToken) string {
	if cli.TenantOverride != "" {
		return cli.TenantOverride
	}

	return token.TenantID()
}

type VersionCmd struct{}

func (cmd *VersionCmd) Run(c *cc.CommonCtx) error {
	fmt.Fprintf(c.OutWriter, "%s - %s (%s)\n",
		x.AppName,
		version.GetInfo().String(),
		x.AppVersionTag,
	)
	return nil
}
