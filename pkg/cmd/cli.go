package cmd

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/version"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

type CLI struct {
	Globals

	Authorizer AuthorizerCmd  `cmd:"" aliases:"a" help:"authorizer commands"`
	Tenant     TenantCmd      `cmd:"" aliases:"t" help:"tenant commands"`
	Directory  DirectoryCmd   `cmd:"" aliases:"d" help:"directory commands"`
	Developer  DeveloperCmd   `cmd:"" aliases:"x" help:"developer commands"`
	User       UserCmd        `cmd:"" aliases:"u" help:"user commands"`
	Login      user.LoginCmd  `cmd:"" help:"login"`
	Logout     user.LogoutCmd `cmd:"" help:"logout"`
	Config     ConfigCmd      `cmd:"" aliases:"c" help:"configuration commands"`
	Version    VersionCmd     `cmd:"" help:"version information"`
	Help       HelpCmd        `cmd:"" hidden:"" default:"1"`
}

func (cmd *CLI) BeforeApply(c *cc.CommonCtx, g *Globals) error {
	if err := c.SetEnv(g.Environment); err != nil {
		return errors.Wrapf(err, "set environment [%s]", g.Environment)
	}
	if g.TenantOverride != "" {
		c.Override(x.TenantIDOverride, g.TenantOverride)
	}
	if g.AuthorizerOverride != "" {
		c.Override(x.AuthorizerOverride, g.AuthorizerOverride)
	}
	return nil
}

func (cmd *CLI) Run() error {
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

type HelpCmd struct{}

func (cmd *HelpCmd) Run(c *cc.CommonCtx) error {
	fmt.Fprintf(c.OutWriter, "No arguments provided, run \"aserto --help\" for more information.\n")
	return nil
}
