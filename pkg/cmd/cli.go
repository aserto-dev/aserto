package cmd

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
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
	ControlPlane ControlPlaneCmd `cmd:"" aliases:"c" help:"control plane commands"`
	User         UserCmd         `cmd:"" aliases:"u" help:"user commands"`
	Login        user.LoginCmd   `cmd:"" help:"login"`
	Logout       user.LogoutCmd  `cmd:"" help:"logout"`
	Config       ConfigCmd       `cmd:"" aliases:"c" help:"configuration commands"`
	Version      VersionCmd      `cmd:"" help:"version information"`

	// ConfigFileMapper implements the `type:"conf"` tag.
	Cfg            string `name:"config" short:"c" type:"conf" env:"ASERTO_ENV" help:"name or path of configuration file"`
	Verbosity      int    `short:"v" type:"counter" help:"Use to increase output verbosity."`
	TenantOverride string `name:"tenant" env:"ASERTO_TENANT_ID" help:"tenant id override"`
}

type ServiceOptions interface {
	Override(svc x.Service, overrides clients.Overrides)
	RequireToken()
}

func (cli *CLI) ConfigOverrider(conf *config.Config) {
	if cli.TenantOverride != "" {
		conf.TenantID = cli.TenantOverride
	}
}

type VersionCmd struct{}

func (cmd *VersionCmd) Run(c *cc.CommonCtx) error {
	fmt.Fprintf(c.UI.Output(), "%s - %s (%s)\n",
		x.AppName,
		version.GetInfo().String(),
		x.AppVersionTag,
	)
	return nil
}
