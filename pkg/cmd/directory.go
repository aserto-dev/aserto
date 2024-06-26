package cmd

import (
	"errors"

	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	errs "github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/topaz/pkg/cli/cmd/directory"
	"github.com/go-http-utils/headers"
	"github.com/samber/lo"
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

	isTopazConfig := !cc.IsAsertoAccount(cfg.ConfigName)

	if isTopazConfig {
		err = setServicesConfig(cfg, c.Config.Active.ConfigFile)
		if err != nil {
			return err
		}
	}

	token, err := getTenantToken(cfg.Auth)
	if err != nil && !errors.Is(err, errs.NeedLoginErr) {
		return err
	}

	dirConfig := topazClients.DirectoryConfig{
		Host:     cfg.Services.DirectoryReaderService.Address,
		APIKey:   cfg.Services.DirectoryReaderService.APIKey,
		Token:    lo.Ternary(isTopazConfig, "", token.Access), // only send access token to hosted services.
		Insecure: cfg.Services.DirectoryReaderService.Insecure,
		TenantID: lo.Ternary(isTopazConfig, "", cfg.TenantID),
	}

	if !isTopazConfig {
		// only send access token to hosted services.
		c.Context = metadata.AppendToOutgoingContext(c.Context, string(headers.Authorization), BearerToken+token.Access)
	}

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
	cmd.DirectoryCmd.Import.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Export.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Backup.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Restore.DirectoryConfig = dirConfig
	cmd.DirectoryCmd.Test.Exec.DirectoryConfig = dirConfig

	return nil
}

func (cmd *DirectoryCmd) Run(c *cc.CommonCtx) error {
	return nil
}
