package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
)

type ConfigCmd struct {
}

func (cmd *ConfigCmd) BeforeApply(c *cc.CommonCtx) error {
	return c.VerifyLoggedIn()
}

func (cmd *ConfigCmd) Run(c *cc.CommonCtx) error {
	return nil
}
