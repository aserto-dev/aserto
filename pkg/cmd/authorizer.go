package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/authorizer"
)

type AuthorizerCmd struct {
	EvalDecision authorizer.EvalDecisionCmd `cmd:"" help:"evaluate policy decision" group:"authorizer"`
	DecisionTree authorizer.DecisionTreeCmd `cmd:"" help:"get decision tree" group:"authorizer"`
	ExecQuery    authorizer.ExecQueryCmd    `cmd:"" help:"execute query" group:"authorizer"`
}

func (cmd *AuthorizerCmd) BeforeApply(c *cc.CommonCtx) error {
	return c.VerifyLoggedIn()
}

func (cmd *AuthorizerCmd) Run(c *cc.CommonCtx) error {
	return nil
}
