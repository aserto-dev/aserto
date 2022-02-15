package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/authorizer"
	"github.com/aserto-dev/aserto/pkg/x"
)

type AuthorizerCmd struct {
	EvalDecision authorizer.EvalDecisionCmd `cmd:"" help:"evaluate policy decision" group:"authorizer"`
	DecisionTree authorizer.DecisionTreeCmd `cmd:"" help:"get decision tree" group:"authorizer"`
	ExecQuery    authorizer.ExecQueryCmd    `cmd:"" help:"execute query" group:"authorizer"`

	AuthorizerOverrides AuthorizerOptions `embed:"" envprefix:"ASERTO_AUTHORIZER_"`
}

func (cmd *AuthorizerCmd) AfterApply(co ConnectionOverrides) error {
	co.Override(x.AuthorizerService, &cmd.AuthorizerOverrides)

	return nil
}

func (cmd *AuthorizerCmd) Run(c *cc.CommonCtx) error {
	return nil
}
