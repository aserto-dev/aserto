package main

import (
	"io/ioutil"
	"log"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cmd"
	"github.com/aserto-dev/aserto/pkg/x"
)

func main() {
	log.SetOutput(ioutil.Discard)

	cli := cmd.CLI{
		Globals: cmd.Globals{},
	}

	c := cc.New()

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
			NoExpandSubcommands: true,
		}),
		kong.Bind(c, &cli.Globals),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
