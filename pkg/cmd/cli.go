package cmd

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/version"
	"github.com/aserto-dev/aserto/pkg/x"
)

type CLI struct {
	Authorizer         AuthorizerCmd     `cmd:"" aliases:"a" help:"authorizer commands"`
	Tenant             TenantCmd         `cmd:"" aliases:"t" help:"tenant commands"`
	Directory          DirectoryCmd      `cmd:"" aliases:"d" help:"directory commands"`
	DecisionLogs       DecisionLogsCmd   `cmd:"" aliases:"l" help:"decision logs commands"`
	Developer          DeveloperCmd      `cmd:"" aliases:"x" help:"developer commands"`
	User               UserCmd           `cmd:"" aliases:"u" help:"user commands"`
	Login              user.LoginCmd     `cmd:"" help:"login"`
	Logout             user.LogoutCmd    `cmd:"" help:"logout"`
	Config             ConfigCmd         `cmd:"" aliases:"c" help:"configuration commands"`
	Version            VersionCmd        `cmd:"" help:"version information"`
	Verbose            bool              `name:"verbose" help:"verbose output"`
	Insecure           insecure          `name:"insecure" help:"skip TLS verification"`
	AuthorizerOverride authorizerAddress `name:"authorizer" env:"ASERTO_AUTHORIZER" help:"authorizer override"`
	TenantOverride     tenantID          `name:"tenant" env:"ASERTO_TENANT_ID" help:"tenant id override"`
	EnvOverride        environment       `name:"env" default:"${defaultEnv}" env:"ASERTO_ENV" hidden:"" help:"environment override"`
	Debug              debug             `name:"debug" env:"ASERTO_DEBUG" help:"enable debug logging"`
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

// An option to disable TLS verification.
type insecure bool

func (ins insecure) AfterApply(c *cc.CommonCtx) error {
	if ins {
		c.Insecure = bool(ins)
	}

	return nil
}

// An option to explicitly set the address of the authorizer service.
type authorizerAddress string

func (address authorizerAddress) AfterApply(c *cc.CommonCtx) error {
	if address != "" {
		c.Override(x.AuthorizerOverride, string(address))
	}
	return nil
}

// An option to explicitly set the tenant ID.
type tenantID string

func (tenantID tenantID) AfterApply(c *cc.CommonCtx) error {
	if tenantID != "" {
		c.Override(x.TenantIDOverride, string(tenantID))
	}
	return nil
}

// An option to set an Aserto environment (e.g. "prod", "eng").
type environment string

func (env environment) AfterApply(c *cc.CommonCtx, ctx *kong.Context) error {
	err := c.SetEnv(string(env))
	ctx.FatalIfErrorf(err, "set environment [%s]", env)
	return err
}

// A flag to emit debug-level output.
type debug bool

func (dbg debug) AfterApply(c *cc.CommonCtx) error {
	if dbg {
		c.SetLogger(os.Stderr)
	}
	return nil
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
