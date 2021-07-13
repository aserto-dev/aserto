package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cmd"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

func main() {
	c := cc.New()

	cli := cmd.CLI{}
	ctx := kong.Parse(&cli,
		kong.Name(x.AppName),
		kong.Description(x.AppDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary:        false,
			Summary:             false,
			Compact:             true,
			Tree:                false,
			FlagsLast:           true,
			Indenter:            kong.SpaceIndenter,
			NoExpandSubcommands: false,
		}),
		kong.Vars{"defaultEnv": x.DefaultEnv},
		kong.Bind(&cli),
	)

	if cli.Debug {
		c.SetLogger(os.Stderr)
	}

	if err := c.SetEnv(cli.EnvOverride); err != nil {
		ctx.FatalIfErrorf(errors.Wrapf(err, "set environment [%s]", cli.EnvOverride))
	}

	if cli.TenantOverride != "" {
		c.Override(x.TenantIDOverride, cli.TenantOverride)
	}

	if cli.AuthorizerOverride != "" {
		c.Override(x.AuthorizerOverride, cli.AuthorizerOverride)
	}

	if cli.IsLoginRequired() {
		if err := c.VerifyLoggedIn(); err != nil {
			fmt.Fprintln(c.ErrWriter, err.Error())
			os.Exit(1)
		}
	}

	err := ctx.Run(c)
	ctx.FatalIfErrorf(err)
}
