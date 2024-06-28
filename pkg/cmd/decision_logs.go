package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/errors"
	dl "github.com/aserto-dev/aserto/pkg/handlers/decision_logs"
	"github.com/aserto-dev/aserto/pkg/x"
)

type DecisionLogsCmd struct {
	List    dl.ListCmd        `cmd:"list" help:"list available decision log files"`
	Get     dl.GetCmd         `cmd:"get" help:"download one or more decision log files"`
	Stream  dl.StreamCmd      `cmd:"stream" help:"stream decision log events to stdout"`
	SvcOpts ConnectionOptions `embed:"" envprefix:"ASERTO_DECISION_LOGS_"`
}

func (cmd *DecisionLogsCmd) BeforeApply(context *kong.Context) error {
	cfg, err := getConfig(context)
	if err != nil {
		return err
	}
	if !cc.IsAsertoAccount(cfg.ConfigName) && cfg.TenantID == "" {
		return errors.ErrDecisionLogsCmd
	}
	return nil
}

func (cmd *DecisionLogsCmd) AfterApply(so ServiceOptions) error {
	so.Override(x.DecisionLogsService, &cmd.SvcOpts)

	return nil
}

func (cmd *DecisionLogsCmd) Run(c *cc.CommonCtx) error {
	return nil
}
