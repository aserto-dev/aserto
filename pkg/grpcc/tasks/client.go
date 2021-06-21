package tasks

import (
	"context"

	"github.com/aserto-dev/aserto/pkg/grpcc"
	tasksmgr "github.com/aserto-dev/proto/aserto/task/manager"
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

// TasksManagerClient -- return tasks manager client.
func (c *Client) TasksManagerClient() tasksmgr.ManagerClient {
	return tasksmgr.NewManagerClient(c.conn)
}
