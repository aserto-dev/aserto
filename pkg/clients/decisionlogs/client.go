package decisionlogs

import (
	"context"

	"github.com/aserto-dev/go-aserto/client"
	decision_logs "github.com/aserto-dev/go-decision-logs/aserto/decision-logs/v2"
	"google.golang.org/grpc"
)

// type Config struct {
// 	Host     string `flag:"host" short:"H" default:"${directory_svc}" env:"TOPAZ_DIRECTORY_SVC" help:"directory service address"`
// 	APIKey   string `flag:"api-key" short:"k" default:"${directory_key}" env:"TOPAZ_DIRECTORY_KEY" help:"directory API key"`
// 	Token    string `flag:"token" default:"${directory_token}" env:"TOPAZ_DIRECTORY_TOKEN" help:"directory OAuth2.0 token" hidden:""`
// 	Insecure bool   `flag:"insecure" short:"i" default:"${insecure}" env:"TOPAZ_INSECURE" help:"skip TLS verification"`
// 	TenantID string `flag:"tenant-id" help:"" default:"${tenant_id}" env:"ASERTO_TENANT_ID" `
// }

type Client struct {
	conn *grpc.ClientConn
	decision_logs.DecisionLogsClient
}

func NewClient(ctx context.Context, options ...client.ConnectionOption) (*Client, error) {
	conn, err := client.NewConnection(ctx, options...)
	if err != nil {
		return nil, err
	}

	return New(conn), nil
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		conn:               conn,
		DecisionLogsClient: decision_logs.NewDecisionLogsClient(conn),
	}
}
