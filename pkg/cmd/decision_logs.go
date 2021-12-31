package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/decision_logs"
)

type DecisionLogsCmd struct {
	List      decision_logs.ListCmd      `cmd:"" help:"list available decision log files" group:"decision-logs"`
	Get       decision_logs.GetCmd       `cmd:"" help:"download one or more decision log files" group:"decision-logs"`
	ListUsers decision_logs.ListUsersCmd `cmd:"" help:"list available user data files" group:"decision-logs"`
	GetUser   decision_logs.GetUserCmd   `cmd:"" help:"download one or more user data files" group:"decision-logs"`
	APIKey    string                     `optional:"" help:"authenticate with an API key"`
}

func (cmd *DecisionLogsCmd) AfterApply(kctx *kong.Context, cli *CLI) error {
	if cmd.APIKey == "" {
		cli.RequireLogin()
	}
	kctx.Bind(decision_logs.APIKey(cmd.APIKey))

	return nil
}

func (cmd *DecisionLogsCmd) Run(c *cc.CommonCtx) error {
	return nil
}
