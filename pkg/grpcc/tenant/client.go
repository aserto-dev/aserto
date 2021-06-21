package tenant

import (
	"context"

	"github.com/aserto-dev/aserto/pkg/grpcc"

	"github.com/aserto-dev/proto/aserto/common/info"
	"github.com/aserto-dev/proto/aserto/tenant/account"
	"github.com/aserto-dev/proto/aserto/tenant/connection"
	"github.com/aserto-dev/proto/aserto/tenant/onboarding"
	"github.com/aserto-dev/proto/aserto/tenant/policy"
	"github.com/aserto-dev/proto/aserto/tenant/profile"
	"github.com/aserto-dev/proto/aserto/tenant/provider"
	"github.com/aserto-dev/proto/aserto/tenant/scc"
	"github.com/aserto-dev/proto/aserto/tenant/system"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client tenant gRPC connection
type Client struct {
	conn *grpc.ClientConn
	addr string
}

func Connection(ctx context.Context, addr string, creds credentials.PerRPCCredentials) (*Client, error) {
	gconn, err := grpcc.NewClient(ctx, addr, creds)
	if err != nil {
		return nil, errors.Wrap(err, "create grpc client failed")
	}

	return &Client{
		conn: gconn.Conn,
		addr: addr,
	}, err
}

func (c *Client) AccountClient() account.AccountClient {
	return account.NewAccountClient(c.conn)
}

func (c *Client) ConnectionManagerClient() connection.ConnectionClient {
	return connection.NewConnectionClient(c.conn)
}

func (c *Client) OnboardingClient() onboarding.OnboardingClient {
	return onboarding.NewOnboardingClient(c.conn)
}

func (c *Client) PolicyClient() policy.PolicyClient {
	return policy.NewPolicyClient(c.conn)
}

func (c *Client) ProfileClient() profile.ProfileClient {
	return profile.NewProfileClient(c.conn)
}

func (c *Client) ProviderClient() provider.ProviderClient {
	return provider.NewProviderClient(c.conn)
}

func (c *Client) SCCClient() scc.SourceCodeCtlClient {
	return scc.NewSourceCodeCtlClient(c.conn)
}

func (c *Client) InfoClient() info.InfoClient {
	return info.NewInfoClient(c.conn)
}

func (c *Client) SystemClient() system.SystemClient {
	return system.NewSystemClient(c.conn)
}
