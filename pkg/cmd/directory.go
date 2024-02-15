package cmd

import (
	"os"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/topaz/pkg/cli/cmd/directory"
	"github.com/go-http-utils/headers"
	"google.golang.org/grpc/metadata"

	aErr "github.com/aserto-dev/aserto/pkg/cc/errors"
	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazClients "github.com/aserto-dev/topaz/pkg/cli/clients"
	topaz "github.com/aserto-dev/topaz/pkg/cli/cmd"
)

type DirectoryCmd struct {
	GetManifest     topaz.GetManifestCmd         `cmd:"" help:"get manifest" group:"directory"`
	SetManifest     topaz.SetManifestCmd         `cmd:"" help:"set manifest" group:"directory"`
	DeleteManifest  topaz.DeleteManifestCmd      `cmd:"" help:"delete manifest" group:"directory"`
	GetObject       directory.GetObjectCmd       `cmd:"" help:"get object" group:"directory"`
	SetObject       directory.SetObjectCmd       `cmd:"" help:"set object" group:"directory"`
	DeleteObject    directory.DeleteObjectCmd    `cmd:"" help:"delete object" group:"directory"`
	ListObjects     directory.ListObjectsCmd     `cmd:"" help:"list objects" group:"directory"`
	GetRelation     directory.GetRelationCmd     `cmd:"" help:"get relation" group:"directory"`
	SetRelation     directory.SetRelationCmd     `cmd:"" help:"set relation" group:"directory"`
	DeleteRelation  directory.DeleteRelationCmd  `cmd:"" help:"delete relation" group:"directory"`
	ListRelations   directory.ListRelationsCmd   `cmd:"" help:"list relations" group:"directory"`
	CheckRelation   directory.CheckRelationCmd   `cmd:"" help:"check relation" group:"directory"`
	CheckPermission directory.CheckPermissionCmd `cmd:"" help:"check permission" group:"directory"`
	GetGraph        directory.GetGraphCmd        `cmd:"" help:"get relation graph" group:"directory"`
}

func (cmd *DirectoryCmd) AfterApply(c *topazCC.CommonCtx) error {
	cfgPath, err := config.GetSymlinkConfigPath()
	if err != nil {
		return err
	}

	if !filex.FileExists(cfgPath) {
		return aErr.NeedLoginErr
	}

	cfg, err := config.NewConfig(config.Path(cfgPath))
	if err != nil {
		return err
	}

	for _, ctxs := range cfg.Context.Contexts {
		if cfg.Context.ActiveContext == ctxs.Name {
			if ctxs.TopazConfigFile != "" {
				err = setServicesConfig(cfg, ctxs.TopazConfigFile)
				if err != nil {
					return err
				}
			}
			err = os.Setenv(topazClients.EnvTopazDirectorySvc, cfg.Services.DirectoryReaderService.Address)
			if err != nil {
				return err
			}

			tenantToken, err := getTenantTokenDetails(ctxs.TenantID, cfg.Auth)
			if err != nil {
				return err
			}

			dirConfig := topazClients.Config{
				Host:     cfg.Services.DirectoryReaderService.Address,
				APIKey:   tenantToken.DirectoryWriteKey,
				Insecure: cfg.Services.DirectoryReaderService.Insecure,
				TenantID: ctxs.TenantID,
			}

			c.Context = metadata.AppendToOutgoingContext(c.Context, string(headers.Authorization), "Basic "+tenantToken.DirectoryWriteKey)

			cmd.GetManifest.Config = dirConfig
			cmd.SetManifest.Config = dirConfig
			cmd.DeleteManifest.Config = dirConfig
			cmd.GetObject.Config = dirConfig
			cmd.SetObject.Config = dirConfig
			cmd.DeleteObject.Config = dirConfig
			cmd.ListObjects.Config = dirConfig
			cmd.GetRelation.Config = dirConfig
			cmd.SetRelation.Config = dirConfig
			cmd.DeleteRelation.Config = dirConfig
			cmd.ListRelations.Config = dirConfig
			cmd.CheckRelation.Config = dirConfig
			cmd.CheckPermission.Config = dirConfig
			cmd.GetGraph.Config = dirConfig

		}
	}
	return nil
}

func (cmd *DirectoryCmd) Run(c *cc.CommonCtx) error {
	return nil
}
