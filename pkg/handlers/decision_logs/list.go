package decision_logs //nolint // prefer standardizing name over removing _

import (
	"context"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"
	"google.golang.org/protobuf/proto"
)

type ListCmd struct {
}

func (cmd ListCmd) Run(c *cc.CommonCtx, apiKey APIKey) error {
	ctx := c.Context
	cli, err := newClient(c, apiKey)
	if err != nil {
		return err
	}
	results, err := listDecisionLogs(ctx, cli)
	if err != nil {
		return nil
	}

	return jsonx.OutputJSONPBArray(c.OutWriter, results)
}

func listDecisionLogs(ctx context.Context, cli dl.DecisionLogsClient) ([]proto.Message, error) {
	next := func(ctx context.Context, token string) (*dl.ListDecisionLogsResponse, error) {
		return cli.ListDecisionLogs(ctx, &dl.ListDecisionLogsRequest{
			Page: &api.PaginationRequest{
				Token: token,
			},
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
