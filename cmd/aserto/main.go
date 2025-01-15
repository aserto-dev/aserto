package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cmd"
	"github.com/aserto-dev/aserto/pkg/cmd/conf"
	"github.com/aserto-dev/aserto/pkg/x"
	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topaz "github.com/aserto-dev/topaz/pkg/cli/cmd/common"
	"github.com/aserto-dev/topaz/pkg/cli/fflag"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
)

const (
	rcOK  int = 0
	rcErr int = 1
)

var (
	tmpConfig  *config.Config
	configOnce sync.Once
)

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "--help")
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		sc, ok := sig.(syscall.Signal)
		if !ok {
			sc = 0
		}

		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "%s request received, canceling...\n", sig.String())
		fmt.Fprintln(os.Stderr)

		cancel()

		os.Exit(128 + int(sc))
	}()

	os.Exit(run(ctx))
}

func run(ctx context.Context) (exitCode int) {
	fflag.Init()

	cwd, err := os.Getwd()
	if err != nil {
		return exitErr(errors.Wrap(err, "failed to determine current working directory"))
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return exitErr(errors.Wrap(err, "failed to determine user home directory"))
	}

	configDir := filepath.Join(home, ".config", x.AppName)

	serviceOptions := clients.NewServiceOptions()

	cliConfigFile := filepath.Join(topazCC.GetTopazDir(), topaz.CLIConfigurationFile)

	topazCtx, err := topazCC.NewCommonContext(ctx, true, cliConfigFile)
	if err != nil {
		return exitErr(err)
	}

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
		kong.Bind(topazCtx),
		kong.NamedMapper("conf", conf.ConfigFileMapper(configDir)),
		kong.BindTo(serviceOptions, (*cmd.ServiceOptions)(nil)),
		kong.Vars{
			"topaz_dir":          topazCC.GetTopazDir(),
			"topaz_certs_dir":    topazCC.GetTopazCertsDir(),
			"topaz_cfg_dir":      topazCC.GetTopazCfgDir(),
			"topaz_db_dir":       topazCC.GetTopazDataDir(),
			"container_registry": topazCC.ContainerRegistry(),
			"container_image":    topazCC.ContainerImage(),
			"container_tag":      topazCC.ContainerTag(),
			"container_platform": topazCC.ContainerPlatform(),
			"container_name":     topazCC.ContainerName(topazCtx.Config.Active.ConfigFile),
			"directory_svc":      topazCC.DirectorySvc(),
			"directory_key":      topazCC.DirectoryKey(),
			"directory_token":    topazCC.DirectoryToken(),
			"authorizer_svc":     topazCC.AuthorizerSvc(),
			"authorizer_key":     topazCC.AuthorizerKey(),
			"authorizer_token":   topazCC.AuthorizerToken(),
			"plaintext":          strconv.FormatBool(topazCC.Plaintext()),
			"tenant_id":          topazCC.TenantID(),
			"insecure":           strconv.FormatBool(topazCC.Insecure()),
			"no_check":           strconv.FormatBool(topazCC.NoCheck()),
			"no_color":           strconv.FormatBool(topazCC.NoColor()),
			"cwd":                cwd,
			"timeout":            topazCC.Timeout().String(),
		},
	)

	configPath := config.DefaultConfigFilePath
	if cli.Cfg != "" {
		configPath = cli.Cfg
	}

	c, err := cc.NewCommonCtx(
		topazCtx,
		config.Path(configPath),
		cli.ConfigOverrider,
		serviceOptions.ConfigOverrider,
	)
	if err != nil {
		return exitErr(err)
	}

	if err := kongCtx.Run(c); err != nil {
		return exitErr(err)
	}

	return rcOK
}

func exitErr(err error) int {
	fmt.Fprintln(os.Stderr, err.Error())
	return rcErr
}

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

func loadConfig(kongCtx *kong.Context) (*config.Config, error) {
	allFlags := kongCtx.Flags()
	for _, f := range allFlags {
		if f.Name == "config" {
			configPath := kongCtx.FlagValue(f).(string)
			return config.NewConfig(config.Path(configPath))
		}
	}

	return nil, nil
}
