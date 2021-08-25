package authorizer

import (
	"context"

	"github.com/aserto-dev/aserto/pkg/grpcc"

	authz "github.com/aserto-dev/go-grpc-authz/aserto/authorizer/authorizer/v1"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"
	policy "github.com/aserto-dev/go-grpc/aserto/authorizer/policy/v1"
	system "github.com/aserto-dev/go-grpc/aserto/authorizer/system/v1"
	info "github.com/aserto-dev/go-grpc/aserto/common/info/v1"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client gRPC connection
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

// AuthorizerClient -- return authorizer client.
func (c *Client) AuthorizerClient() authz.AuthorizerClient {
	return authz.NewAuthorizerClient(c.conn)
}

// DirectoryClient -- return directory client.
func (c *Client) DirectoryClient() dir.DirectoryClient {
	return dir.NewDirectoryClient(c.conn)
}

// PolicyClient -- return policy client.
func (c *Client) PolicyClient() policy.PolicyClient {
	return policy.NewPolicyClient(c.conn)
}

// InfoClient -- return information client.
func (c *Client) InfoClient() info.InfoClient {
	return info.NewInfoClient(c.conn)
}

// SystemClient -- return information client.
func (c *Client) SystemClient() system.SystemClient {
	return system.NewSystemClient(c.conn)
}
