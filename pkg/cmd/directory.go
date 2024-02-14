package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	topaz "github.com/aserto-dev/topaz/pkg/cli/cmd"
	"github.com/aserto-dev/topaz/pkg/cli/cmd/directory"
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

func (cmd *DirectoryCmd) AfterApply(so ServiceOptions) error {
	return nil
}

func (cmd *DirectoryCmd) Run(c *cc.CommonCtx) error {
	return nil
}
