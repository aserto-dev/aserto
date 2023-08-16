package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/directory"
)

type DirectoryCmd struct {
	GetManifest     directory.GetManifestCmd     `cmd:"" help:"get manifest" group:"directory"`
	SetManifest     directory.SetManifestCmd     `cmd:"" help:"set manifest" group:"directory"`
	DeleteManifest  directory.DeleteManifestCmd  `cmd:"" help:"delete manifest" group:"directory"`
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
	Load            directory.LoadCmd            `cmd:"" help:"load data into the directory" group:"directory"`
	// AuthorizerOverrides AuthorizerOptions         `embed:"" envprefix:"ASERTO_AUTHORIZER_"`
}

func (cmd *DirectoryCmd) AfterApply(so ServiceOptions) error {
	// so.Override(x.AuthorizerService, &cmd.AuthorizerOverrides)
	return nil
}

func (cmd *DirectoryCmd) Run(c *cc.CommonCtx) error {
	return nil
}
