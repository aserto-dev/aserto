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
	dsClient "github.com/aserto-dev/topaz/pkg/cli/clients/directory"
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
	if err != nil && !errors.Is(err, errs.ErrNeedLogin) {
		return err
	}

	dirConfig := dsClient.Config{
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

	cmd.DirectoryCmd.Get.Manifest.Config = dirConfig
	cmd.DirectoryCmd.Set.Manifest.Config = dirConfig
	cmd.DirectoryCmd.Delete.Manifest.Config = dirConfig
	cmd.DirectoryCmd.Get.Object.Config = dirConfig
	cmd.DirectoryCmd.Set.Object.Config = dirConfig
	cmd.DirectoryCmd.Delete.Object.Config = dirConfig
	cmd.DirectoryCmd.List.Objects.Config = dirConfig
	cmd.DirectoryCmd.Get.Relation.Config = dirConfig
	cmd.DirectoryCmd.Set.Relation.Config = dirConfig
	cmd.DirectoryCmd.Delete.Relation.Config = dirConfig
	cmd.DirectoryCmd.List.Relations.Config = dirConfig
	cmd.DirectoryCmd.Check.Config = dirConfig
	cmd.DirectoryCmd.Search.Config = dirConfig
	cmd.DirectoryCmd.Import.Config = dirConfig
	cmd.DirectoryCmd.Export.Config = dirConfig
	cmd.DirectoryCmd.Backup.Config = dirConfig
	cmd.DirectoryCmd.Restore.Config = dirConfig
	cmd.DirectoryCmd.Test.Exec.Config = dirConfig
	cmd.DirectoryCmd.Stats.Config = dirConfig

	return nil
}

func (cmd *DirectoryCmd) Run(c *cc.CommonCtx) error {
	return nil
}
