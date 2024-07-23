package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cmd"
	"github.com/aserto-dev/aserto/pkg/cmd/conf"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/aserto-dev/topaz/pkg/cli/fflag"
	"github.com/pkg/errors"

	"github.com/aserto-dev/go-aserto/client"

	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topaz "github.com/aserto-dev/topaz/pkg/cli/cmd/common"
)

func main() {
	fflag.Init()

	home, err := os.UserHomeDir()
	if err != nil {
		panic(errors.Wrap(err, "failed to determine user home directory"))
	}

	configDir := filepath.Join(home, ".config", x.AppName)

	serviceOptions := clients.NewServiceOptions()

	cli := cmd.CLI{}

	cliConfigFile := filepath.Join(topazCC.GetTopazDir(), topaz.CLIConfigurationFile)
	topazCtx, err := topazCC.NewCommonContext(true, cliConfigFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	containerVersion := topazCC.ContainerTag()
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if strings.Contains(dep.Path, "github.com/aserto-dev/topaz") {
				containerVersion = strings.TrimPrefix(dep.Version, "v")
			}
		}
	}

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
		kong.Bind(topazCtx),
		kong.NamedMapper("conf", conf.ConfigFileMapper(configDir)), // attach to tag `type:"conf"`
		kong.BindTo(serviceOptions, (*cmd.ServiceOptions)(nil)),
		kong.Vars{
			"topaz_dir":          topazCC.GetTopazDir(),
			"topaz_certs_dir":    topazCC.GetTopazCertsDir(),
			"topaz_cfg_dir":      topazCC.GetTopazCfgDir(),
			"topaz_db_dir":       topazCC.GetTopazDataDir(),
			"container_registry": topazCC.ContainerRegistry(),
			"container_image":    topazCC.ContainerImage(),
			"container_tag":      containerVersion,
			"container_platform": topazCC.ContainerPlatform(),
			"container_name":     topazCC.ContainerName(topazCtx.Config.Active.ConfigFile),
			"directory_svc":      topazCC.DirectorySvc(),
			"directory_key":      topazCC.DirectoryKey(),
			"directory_token":    topazCC.DirectoryToken(),
			"authorizer_svc":     topazCC.AuthorizerSvc(),
			"authorizer_key":     topazCC.AuthorizerKey(),
			"authorizer_token":   topazCC.AuthorizerToken(),
			"tenant_id":          topazCC.TenantID(),
			"insecure":           strconv.FormatBool(topazCC.Insecure()),
			"no_check":           strconv.FormatBool(topazCC.NoCheck()),
			"no_color":           strconv.FormatBool(topazCC.NoColor()),
		},
	)
	configPath := config.DefaultConfigFilePath
	if cli.Cfg != "" {
		configPath = cli.Cfg
	}

	ctx, err := cc.NewCommonCtx(
		topazCtx,
		config.Path(configPath),
		cli.ConfigOverrider,
		serviceOptions.ConfigOverrider,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	topazCtx.Context = client.SetTenantContext(topazCtx.Context, ctx.TenantID())
	ctx.CommonCtx = topazCtx

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
		case "ASERTO_AUTHORIZER_":
			svcOptions = &tmpConfig.Services.AuthorizerService
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
		case "authorizer":
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

	return nil, nil
}
