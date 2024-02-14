package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	topazAuthz "github.com/aserto-dev/topaz/pkg/cli/cmd/authorizer"
)

type AuthorizerCmd struct {
	EvalDecision topazAuthz.EvalDecisionCmd `cmd:"" help:"evaluate policy decision" group:"authorizer"`
	DecisionTree topazAuthz.DecisionTreeCmd `cmd:"" help:"get decision tree" group:"authorizer"`
	ExecQuery    topazAuthz.ExecQueryCmd    `cmd:"" help:"execute query" group:"authorizer"`
	Compile      topazAuthz.CompileCmd      `cmd:"" help:"compile query" group:"authorizer"`
	GetPolicy    topazAuthz.GetPolicyCmd    `cmd:"" help:"get policy" group:"authorizer"`
	ListPolicies topazAuthz.ListPoliciesCmd `cmd:"" help:"list policies" group:"authorizer"`
}

func (cmd *AuthorizerCmd) Run(c *cc.CommonCtx) error {
	return nil
}
