package decision_logs //nolint // prefer standardizing name over removing _

import (
	"os"
	"os/signal"
	"time"

	"github.com/aserto-dev/aserto/pkg/cc"
	dl "github.com/aserto-dev/go-decision-logs/aserto/decision-logs/v2"
	"github.com/aserto-dev/topaz/pkg/cli/jsonx"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type StreamCmd struct {
	PolicyName    string `arg:"" help:"Name of policy to open stream for"`
	InstanceLabel string `arg:"" help:"Label of policy to open stream for"`
	Since         string `optional:"" help:"time to start streaming events from in RFC3339 format"`
}

func (cmd StreamCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.DecisionLogsClient(c.Context)
	if err != nil {
		return err
	}

	var sincePB *timestamppb.Timestamp
	if cmd.Since != "" {
		sinceTime, parseErr := time.Parse(time.RFC3339, cmd.Since)
		if parseErr != nil {
			return parseErr
		}
		sincePB = timestamppb.New(sinceTime)
	}

	stream, err := cli.GetDecisions(c.Context, &dl.GetDecisionsRequest{
		PolicyName:    cmd.PolicyName,
		InstanceLabel: cmd.InstanceLabel,
		Since:         sincePB,
	})
	if err != nil {
		return err
	}

	done := make(chan os.Signal, 1)
	errCh := make(chan error)
	signal.Notify(done, os.Interrupt)

	go func() {
		for {
			resp, errRcv := stream.Recv()
			if errRcv != nil {
				errCh <- errRcv
				return
			}

			errRcv = jsonx.OutputJSON(c.StdOut(), resp.Decision)
			if err != nil {
				errCh <- errRcv
				return
			}
		}
	}()

	select {
	case err = <-errCh:
		return err
	case <-done:
	}

	return nil
}
