package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cmd"
	"github.com/aserto-dev/aserto/pkg/filex"
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
		kong.NamedMapper("conf", configFileMapper(configDir)), // attach to tag `type:"conf"`
		kong.BindTo(serviceOptions, (*cmd.ServiceOptions)(nil)),
	)

	ctx, err := cc.BuildCommonCtx(
		config.Path(cli.Cfg),
		configOverrider(&cli),
		serviceOptions,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := kongCtx.Run(ctx); err != nil {
		kongCtx.FatalIfErrorf(err)
	}
}

func configOverrider(cli *cmd.CLI) config.Overrider {
	return func(conf *config.Config) {
		if cli.TenantOverride != "" {
			conf.TenantID = cli.TenantOverride
		}
	}
}

// configFileMapper is a kong.Mapper that resolves config files.
//
// When applied to a CLI flag, it attempts to find a configuration file that best matches the specified name using
// the following rules:
// 1. If the value is a full or relative path to an existing file, that file is chosen.
// 2. If the value is a file name (without a path separator) with an extension (e.g. "config.yaml") and a file with that
//    name exists in the config directory, that file is chosen.
// 3. If the value is a string without an dot (e.g. "eng") and a file with that name (i.e. "eng.*") exists in the
//    config directory, that file is chosen. If multiple files match, the first one is chosen and a warning is printed
//    to stderr.
type configFileMapper string

func (m configFileMapper) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	if target.Kind() != reflect.String {
		return errors.Errorf(`"conf" type must be applied to a string not %s`, target.Type())
	}

	var path string
	if err := ctx.Scan.PopValueInto("file", &path); err != nil {
		return err
	}

	if path != "-" {
		path = m.find(path)
	}

	target.SetString(path)
	return nil
}

func (m configFileMapper) find(path string) string {
	expanded := kong.ExpandPath(path)
	if filex.FileExists(expanded) {
		return expanded
	}

	if !strings.ContainsRune(path, filepath.Separator) {
		expanded = filepath.Join(string(m), path)
		if filepath.Ext(path) != "" {
			// It's a filename with no path. Look in config directory.
			if filex.FileExists(expanded) {
				return expanded
			}
		} else if matches, err := filepath.Glob(expanded + ".*"); err == nil && len(matches) > 0 {
			if len(matches) > 1 {
				fmt.Fprintf(
					os.Stderr,
					"WARNING: The specified configuration ('%s') matches multiple configuration files: %s. Using '%s'",
					path,
					matches,
					matches[0],
				)
			}
			return matches[0]
		}
	}

	return path
}
