package directory

import "github.com/aserto-dev/aserto/pkg/cc"

type GetObjectCmd struct{}

func (cmd *GetObjectCmd) Run(c *cc.CommonCtx) error {
	return nil
}

type SetObjectCmd struct{}

func (cmd *SetObjectCmd) Run(c *cc.CommonCtx) error {
	return nil
}

type DeleteObjectCmd struct{}

func (cmd *DeleteObjectCmd) Run(c *cc.CommonCtx) error {
	return nil
}

type ListObjectsCmd struct{}

func (cmd *ListObjectsCmd) Run(c *cc.CommonCtx) error {
	return nil
}
