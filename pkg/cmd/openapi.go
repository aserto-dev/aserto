package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/openapi"
)

type OpenAPICmd struct {
	GeneratePolicy openapi.GenerateOpenAPI `cmd:"" group:"openapi" help:"generate an open api policy"`
}

func (cmd *OpenAPICmd) Run(c *cc.CommonCtx) error {
	return nil
}
