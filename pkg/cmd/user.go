package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
)

type UserCmd struct {
}

func (cmd *UserCmd) BeforeApply(so ServiceOptions) error {
	so.RequireToken()
	return nil
}

func (cmd *UserCmd) Run(c *cc.CommonCtx) error {
	return nil
}
