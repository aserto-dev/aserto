package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cmd"
	"github.com/aserto-dev/aserto/pkg/x"
)

func main() {
	factoryBuilder := cc.NewClientFactoryBuilder()

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
		kong.Vars{"defaultEnv": x.DefaultEnv},
		kong.BindTo(factoryBuilder, (*cmd.ConnectionOverrides)(nil)),
	)

	configureLogging(cli.Debug)

	env, err := x.Environment(cli.EnvOverride)
	if err != nil {
		fatal(err)
	}

	token := cc.NewCachedToken(env.Environment)
	tenantID := cli.TenantID(token)

	clientFactory, err := factoryBuilder.ClientFactory(env, tenantID, token)
	if err != nil {
		fatal(err)
	}

	c := cc.NewCommonCtx(env, tenantID, clientFactory, token.Get())
	if err := kongCtx.Run(c); err != nil {
		kongCtx.FatalIfErrorf(err)
	}
}

func configureLogging(debug bool) {
	log.SetOutput(logWriter(debug))
	log.SetPrefix("")
	log.SetFlags(log.LstdFlags)
}

func logWriter(debug bool) io.Writer {
	if debug {
		return os.Stderr
	}

	return io.Discard
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
