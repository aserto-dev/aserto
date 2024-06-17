package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/topaz/pkg/cli/cmd/directory"
	"github.com/go-http-utils/headers"
	"google.golang.org/grpc/metadata"

	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazClients "github.com/aserto-dev/topaz/pkg/cli/clients"
)

type DirectoryCmd struct {
	directory.DirectoryCmd
}

func (cmd *DirectoryCmd) AfterApply(context *kong.Context, c *topazCC.CommonCtx) error {
	cfg, err := getConfig(context)
	if err != nil {
		return err
	}

	if !cc.IsAsertoAccount(cfg.ConfigName) {
		err = setServicesConfig(cfg, c.Config.Active.ConfigFile)
		if err != nil {
			return err
		}
	}
	tenantToken, err := getTenantTokenDetails(cfg.Auth)
	if err != nil {
		return err
	}
	useTenantID := ""
	if tenantToken == "" {
		useTenantID = cfg.TenantID
	}

	dirConfig := topazClients.DirectoryConfig{
		Host:     cfg.Services.DirectoryReaderService.Address,
		APIKey:   cfg.Services.DirectoryReaderService.APIKey,
		Token:    "",
		Insecure: cfg.Services.DirectoryReaderService.Insecure,
		TenantID: useTenantID,
	}

	c.Context = metadata.AppendToOutgoingContext(c.Context, string(headers.Authorization), "Bearer "+tenantToken)

	cmd.DirectoryCmd.Get.Manifest.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Set.Manifest.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Delete.Manifest.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Get.Object.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Set.Object.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Delete.Object.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.List.Objects.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Get.Relation.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Set.Relation.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Delete.Relation.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.List.Relations.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Check.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Search.DirectoryConfig = dirConfig

	return nil
}

func (cmd *DirectoryCmd) Run(c *cc.CommonCtx) error {
	return nil
}
