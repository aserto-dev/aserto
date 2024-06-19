package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/version"
	"github.com/aserto-dev/aserto/pkg/x"

	topazConfig "github.com/aserto-dev/topaz/pkg/cc/config"
	topaz "github.com/aserto-dev/topaz/pkg/cli/cmd/topaz"
)

type CLI struct {
	Login        user.LoginCmd      `cmd:"" help:"login to aserto.com"`
	Logout       user.LogoutCmd     `cmd:"" help:"logout from aserto.com"`
	Start        topaz.StartCmd     `cmd:"" help:"start a topaz instance"`
	Stop         topaz.StopCmd      `cmd:"" help:"stop a topaz instance"`
	Restart      topaz.RestartCmd   `cmd:"" help:"restart topaz instance"`
	Status       topaz.StatusCmd    `cmd:"" help:"status of topaz instance"`
	RunCmd       topaz.RunCmd       `cmd:"" name:"run" help:"run topaz instance in interactive mode"`
	Console      topaz.ConsoleCmd   `cmd:"" help:"launch web console"`
	Config       ConfigCmd          `cmd:"" help:"configuration commands"`
	Directory    DirectoryCmd       `cmd:"" aliases:"ds" help:"directory commands"`
	Authorizer   AuthorizerCmd      `cmd:"" aliases:"az" help:"authorizer commands"`
	DecisionLogs DecisionLogsCmd    `cmd:"" aliases:"dl" help:"decision logs commands"`
	ControlPlane ControlPlaneCmd    `cmd:"" aliases:"cp" help:"control plane commands"`
	Tenant       TenantCmd          `cmd:"" aliases:"tn" help:"tenant commands"`
	User         UserCmd            `cmd:"" help:"user commands"`
	Install      topaz.InstallCmd   `cmd:"" help:"install topaz"`
	Uninstall    topaz.UninstallCmd `cmd:"" help:"uninstall topaz, removes all locally installed artifacts"`
	Update       topaz.UpdateCmd    `cmd:"" help:"update topaz container version"`
	Version      VersionCmd         `cmd:"" help:"version information"`

	// ConfigFileMapper implements the `type:"conf"` tag.
	Cfg            string `name:"config" short:"c" type:"conf" help:"name or path of configuration file"`
	Verbosity      int    `short:"v" type:"counter" help:"Use to increase output verbosity."`
	TenantOverride string `name:"tenant" env:"ASERTO_TENANT_ID" help:"tenant id override"`
}

type ServiceOptions interface {
	Override(svc x.Service, overrides clients.Overrides)
	RequireToken()
}

func (cli *CLI) ConfigOverrider(conf *config.Config) {
	if cli.TenantOverride != "" {
		conf.TenantID = cli.TenantOverride
	}
}

type VersionCmd struct{}

func (cmd *VersionCmd) Run(c *cc.CommonCtx) error {
	fmt.Fprintf(c.UI.Output(), "%s - %s (%s)\n",
		x.AppName,
		version.GetInfo().String(),
		x.AppVersionTag,
	)
	return nil
}

func setServicesConfig(cfg *config.Config, topazConfigFile string) error {
	loader, err := topazConfig.LoadConfiguration(topazConfigFile)
	if err != nil {
		return err
	}
	// Get first API key in configuration.
	for key := range loader.Configuration.Auth.APIKeys {
		cfg.Services.AuthorizerService.APIKey = key
		break
	}
	if authorizerConfig, ok := loader.Configuration.APIConfig.Services["authorizer"]; ok {
		cfg.Services.AuthorizerService.Address = authorizerConfig.GRPC.ListenAddress
		cfg.Services.AuthorizerService.CACertPath = authorizerConfig.GRPC.Certs.TLSCACertPath
		cfg.Services.AuthorizerService.Insecure = true
	}
	if readerConfig, ok := loader.Configuration.APIConfig.Services["reader"]; ok {
		cfg.Services.DirectoryReaderService.Address = readerConfig.GRPC.ListenAddress
		cfg.Services.DirectoryReaderService.CACertPath = readerConfig.GRPC.Certs.TLSCACertPath
		cfg.Services.DirectoryReaderService.Insecure = true
	}
	if writerConfig, ok := loader.Configuration.APIConfig.Services["writer"]; ok {
		cfg.Services.DirectoryWriterService.Address = writerConfig.GRPC.ListenAddress
		cfg.Services.DirectoryWriterService.CACertPath = writerConfig.GRPC.Certs.TLSCACertPath
		cfg.Services.DirectoryWriterService.Insecure = true
	}
	if modelConfig, ok := loader.Configuration.APIConfig.Services["model"]; ok {
		cfg.Services.DirectoryModelService.Address = modelConfig.GRPC.ListenAddress
		cfg.Services.DirectoryModelService.CACertPath = modelConfig.GRPC.Certs.TLSCACertPath
		cfg.Services.DirectoryModelService.Insecure = true
	}
	return nil
}

func getTenantTokenDetails(cfg *config.Auth) (string, error) {
	cachedToken := cc.GetCacheKey(cfg)
	tkn := token.Load(cachedToken)
	authToken, err := tkn.Get()
	if err != nil {
		return "", err
	}
	return authToken.Access, nil
}

func getConfig(context *kong.Context) (*config.Config, error) {
	allFlags := context.Flags()
	for _, f := range allFlags {
		if f.Name != ConfigFlag {
			continue
		}
		configPath := context.FlagValue(f).(string)
		if configPath == "" {
			configPath = config.DefaultConfigFilePath
		}
		cfg, err := config.NewConfig(config.Path(configPath))
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	return nil, errors.ErrConfigNotFound
}
