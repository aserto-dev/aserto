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
	Compile      authorizer.CompileCmd      `cmd:"" help:"compile query" group:"authorizer"`
	GetPolicy    authorizer.GetPolicyCmd    `cmd:"" help:"get policy" group:"authorizer"`
	ListPolicies authorizer.ListPoliciesCmd `cmd:"" help:"list policies" group:"authorizer"`

	AuthorizerOverrides ServiceOverrideOptions `embed:"" envprefix:"ASERTO_SERVICES_AUTHORIZER_"`
}

func (cmd *AuthorizerCmd) AfterApply(so ServiceOptions) error {
	so.Override(x.AuthorizerService, &cmd.AuthorizerOverrides)

	return nil
}

func (cmd *AuthorizerCmd) Run(c *cc.CommonCtx) error {
	return nil
}
