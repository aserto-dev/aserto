package decision_logs //nolint // prefer standardizing name over removing _

import (
	"context"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	dl "github.com/aserto-dev/go-decision-logs/aserto/decision-logs/v2"
	"github.com/aserto-dev/go-decision-logs/aserto/decision-logs/v2/api"
	"google.golang.org/protobuf/proto"
)

type ListCmd struct {
	Policies []string `optional:"" sep:"," help:"Names of policies to list logs for (all if not specified)"`
}

func (cmd ListCmd) Run(c *cc.CommonCtx) error {
	ctx := c.Context
	cli, err := c.DecisionLogsClient()
	if err != nil {
		return err
	}
	results, err := listDecisionLogs(ctx, cli, cmd.Policies)
	if err != nil {
		return err
	}

	return jsonx.OutputJSONPBArray(c.StdOut(), results)
}

func listDecisionLogs(ctx context.Context, cli dl.DecisionLogsClient, policies []string) ([]proto.Message, error) {
	next := func(ctx context.Context, token string) (*dl.ListDecisionLogsResponse, error) {
		return cli.ListDecisionLogs(ctx, &dl.ListDecisionLogsRequest{
			Page: &api.PaginationRequest{
				Token: token,
				Size:  100,
			},
			Policies: policies,
		})
	}

	results := []proto.Message{}

	for resp, err := next(ctx, ""); ; resp, err = next(ctx, resp.Page.NextToken) {
		if err != nil {
			return nil, err
		}
		for _, r := range resp.Results {
			results = append(results, r)
		}
		if resp.Page.NextToken == "" {
			break
		}
	}

	return results, nil
}
