package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cmd"
	"github.com/aserto-dev/aserto/pkg/cmd/conf"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(errors.Wrap(err, "failed to determine user home directory"))
	}

	configDir := filepath.Join(home, ".config", x.AppName)

	serviceOptions := clients.NewServiceOptions()

	cli := cmd.CLI{}
	kongCtx := kong.Parse(&cli,
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
		kong.NamedMapper("conf", conf.ConfigFileMapper(configDir)), // attach to tag `type:"conf"`
		kong.BindTo(serviceOptions, (*cmd.ServiceOptions)(nil)),
	)

	ctx, err := cc.BuildCommonCtx(
		config.Path(cli.Cfg),
		cli.ConfigOverrider,
		serviceOptions.ConfigOverrider,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := kongCtx.Run(ctx); err != nil {
		kongCtx.FatalIfErrorf(err)
	}
}
