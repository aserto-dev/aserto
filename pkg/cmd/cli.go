package cmd

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/version"
	"github.com/aserto-dev/aserto/pkg/x"

	topazConfig "github.com/aserto-dev/topaz/pkg/cc/config"
)

type CLI struct {
	Topaz        TopazCmd        `cmd:"" aliases:"tz" help:"topaz commands"`
	Directory    DirectoryCmd    `cmd:"" aliases:"ds" help:"directory commands"`
	Authorizer   AuthorizerCmd   `cmd:"" aliases:"az" help:"authorizer commands"`
	DecisionLogs DecisionLogsCmd `cmd:"" aliases:"dl" help:"decision logs commands"`
	ControlPlane ControlPlaneCmd `cmd:"" aliases:"cp" help:"control plane commands"`
	Tenant       TenantCmd       `cmd:"" aliases:"tn" help:"tenant commands"`
	Login        user.LoginCmd   `cmd:"" help:"login"`
	Logout       user.LogoutCmd  `cmd:"" help:"logout"`
	Config       ConfigCmd       `cmd:"" help:"configuration commands"`
	Version      VersionCmd      `cmd:"" help:"version information"`

	// ConfigFileMapper implements the `type:"conf"` tag.
	Cfg            string `name:"config" short:"c" type:"conf" env:"ASERTO_ENV" help:"name or path of configuration file"`
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
