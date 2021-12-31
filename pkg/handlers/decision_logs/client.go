package decision_logs //nolint // prefer standardizing name over removing _

import (
	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/client/grpc"
	"github.com/aserto-dev/aserto/pkg/cc"
	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"
)

func newClient(c *cc.CommonCtx, apiKey APIKey) (dl.DecisionLogsClient, error) {
	opts := []aserto.ConnectionOption{
		aserto.WithAddr(c.DecisionLogsService()),
		aserto.WithTenantID(c.TenantID()),
		aserto.WithInsecure(c.Insecure),
	}
	if apiKey != "" {
		opts = append(opts, aserto.WithAPIKeyAuth(string(apiKey)))
	} else {
		opts = append(opts, aserto.WithTokenAuth(c.AccessToken()))
	}

	conn, err := grpc.NewConnection(c.Context, opts...)
	if err != nil {
		return nil, err
	}

	return dl.NewDecisionLogsClient(conn.Conn), nil
}
