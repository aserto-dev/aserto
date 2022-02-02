package decision_logs //nolint // prefer standardizing name over removing _

import (
	"github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto/pkg/cc"
	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"
)

func newClient(c *cc.CommonCtx, apiKey APIKey) (dl.DecisionLogsClient, error) {
	opts := []client.ConnectionOption{
		client.WithAddr(c.DecisionLogsService()),
		client.WithTenantID(c.TenantID()),
		client.WithInsecure(c.Insecure),
	}
	if apiKey != "" {
		opts = append(opts, client.WithAPIKeyAuth(string(apiKey)))
	} else {
		opts = append(opts, client.WithTokenAuth(c.AccessToken()))
	}

	conn, err := client.NewConnection(c.Context, opts...)
	if err != nil {
		return nil, err
	}

	return dl.NewDecisionLogsClient(conn.Conn), nil
}
