package directory

import "github.com/aserto-dev/aserto/pkg/cc"

type GetRelationCmd struct{}

func (cmd *GetRelationCmd) Run(c *cc.CommonCtx) error {
	return nil
}

type SetRelationCmd struct{}

func (cmd *SetRelationCmd) Run(c *cc.CommonCtx) error {
	return nil
}

type DeleteRelationCmd struct{}

func (cmd *DeleteRelationCmd) Run(c *cc.CommonCtx) error {
	return nil
}

type ListRelationsCmd struct{}

func (cmd *ListRelationsCmd) Run(c *cc.CommonCtx) error {
	return nil
}
