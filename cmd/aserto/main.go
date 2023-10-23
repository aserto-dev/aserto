package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"

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
			NoExpandSubcommands: true,
		}),
		kong.Resolvers(ConfigResolver()),
		kong.NamedMapper("conf", conf.ConfigFileMapper(configDir)), // attach to tag `type:"conf"`
		kong.BindTo(serviceOptions, (*cmd.ServiceOptions)(nil)),
	)

	configPath := cli.Login.Cfg
	if configPath == "" {
		configPath = path.Join(configDir, config.ConfigPath)
	}

	ctx, err := cc.BuildCommonCtx(
		config.Path(configPath),
		clients.TenantID(cli.TenantID),
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

var (
	tmpConfig  *config.Config
	configOnce sync.Once
)

// ConfigResolver loads the config file, if present, and populates default values for service connection options like
// address, api-key, insecure, etc.
func ConfigResolver() kong.Resolver {
	var f kong.ResolverFunc = func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (resolved interface{}, err error) {
		configOnce.Do(func() {
			tmpConfig, err = loadConfig(context)
		})

		if err != nil || flag.Tag == nil || flag.Tag.EnvPrefix == "" {
			return resolved, err
		}

		var svcOptions *x.ServiceOptions = nil

		// Only the authorizer and decision logs services have CLI flags to override service options.
		switch flag.Tag.EnvPrefix {
		case "ASERTO_SERVICES_AUTHORIZER_":
			svcOptions = &tmpConfig.Services.AuthorizerService
		case "ASERTO_SERVICES_DIRECTORY_":
			svcOptions = &tmpConfig.Services.DirectoryService
		case "ASERTO_DECISION_LOGS_":
			svcOptions = &tmpConfig.Services.DecisionLogsService
		default:
			return resolved, err
		}

		switch flag.Name {
		case "api-key":
			resolved = svcOptions.APIKey
		case "no-auth":
			flag.Default = strconv.FormatBool(svcOptions.Anonymous)
			resolved = flag.Default
		case "insecure":
			flag.Default = strconv.FormatBool(svcOptions.Insecure)
			resolved = flag.Default
		case "address":
			flag.Default = svcOptions.Address
			resolved = flag.Default
		}

		return resolved, err
	}

	return f
}

func loadConfig(context *kong.Context) (*config.Config, error) {
	allFlags := context.Flags()
	for _, f := range allFlags {
		if f.Name == "config" {
			configPath := context.FlagValue(f).(string)
			return config.NewConfig(config.Path(configPath))
		}
	}

	return config.NewConfig(config.Path(""))
}
