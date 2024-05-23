package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/topaz/pkg/cli/cmd/directory"
	"github.com/go-http-utils/headers"
	"google.golang.org/grpc/metadata"

	aErr "github.com/aserto-dev/aserto/pkg/cc/errors"
	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazClients "github.com/aserto-dev/topaz/pkg/cli/clients"
)

type DirectoryCmd struct {
	GetManifest    directory.GetManifestCmd    `cmd:"" help:"get manifest" group:"directory"`
	SetManifest    directory.SetManifestCmd    `cmd:"" help:"set manifest" group:"directory"`
	DeleteManifest directory.DeleteManifestCmd `cmd:"" help:"delete manifest" group:"directory"`
	GetObject      directory.GetObjectCmd      `cmd:"" help:"get object" group:"directory"`
	SetObject      directory.SetObjectCmd      `cmd:"" help:"set object" group:"directory"`
	DeleteObject   directory.DeleteObjectCmd   `cmd:"" help:"delete object" group:"directory"`
	ListObjects    directory.ListObjectsCmd    `cmd:"" help:"list objects" group:"directory"`
	GetRelation    directory.GetRelationCmd    `cmd:"" help:"get relation" group:"directory"`
	SetRelation    directory.SetRelationCmd    `cmd:"" help:"set relation" group:"directory"`
	DeleteRelation directory.DeleteRelationCmd `cmd:"" help:"delete relation" group:"directory"`
	ListRelations  directory.ListRelationsCmd  `cmd:"" help:"list relations" group:"directory"`
	Check          directory.CheckCmd          `cmd:"" help:"check" group:"directory"`
	Search         directory.SearchCmd         `cmd:"" help:"get relation graph" group:"directory"`
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
		if cfg.Context.ActiveContext != ctxs.Name {
			continue
		}
		if ctxs.TopazConfigFile != "" {
			err = setServicesConfig(cfg, ctxs.TopazConfigFile)
			if err != nil {
				return err
			}
		}

		tenantToken, err := getTenantTokenDetails(cfg.Auth)
		if err != nil {
			return err
		}

		dirConfig := topazClients.DirectoryConfig{
			Host:     cfg.Services.DirectoryReaderService.Address,
			APIKey:   "",
			Token:    tenantToken,
			Insecure: cfg.Services.DirectoryReaderService.Insecure,
			TenantID: ctxs.TenantID,
		}

		c.Context = metadata.AppendToOutgoingContext(c.Context, string(headers.Authorization), "Bearer "+tenantToken)

		cmd.GetManifest.DirectoryConfig = dirConfig
		cmd.SetManifest.DirectoryConfig = dirConfig
		cmd.DeleteManifest.DirectoryConfig = dirConfig
		cmd.GetObject.DirectoryConfig = dirConfig
		cmd.SetObject.DirectoryConfig = dirConfig
		cmd.DeleteObject.DirectoryConfig = dirConfig
		cmd.ListObjects.DirectoryConfig = dirConfig
		cmd.GetRelation.DirectoryConfig = dirConfig
		cmd.SetRelation.DirectoryConfig = dirConfig
		cmd.DeleteRelation.DirectoryConfig = dirConfig
		cmd.ListRelations.DirectoryConfig = dirConfig
		cmd.Check.DirectoryConfig = dirConfig
		cmd.Search.DirectoryConfig = dirConfig

	}
	return nil
}

func (cmd *DirectoryCmd) Run(c *cc.CommonCtx) error {
	return nil
}
