package decision_logs //nolint // prefer standardizing name over removing _

import (
	"os"
	"os/signal"
	"time"

	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StreamCmd struct {
	Policy string `arg:"" help:"ID of policy to open stream for"`
	Since  string `optional:"" help:"UTC time to start streaming events from in RFC3339 format"`
}

func (cmd StreamCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.DecisionLogsClient()
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
		PolicyId: cmd.Policy,
		Since:    sincePB,
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
			}

			errRcv = jsonx.OutputJSON(c.UI.Output(), resp.Decision)
			if err != nil {
				errCh <- errRcv
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
