package directory

import (
	"context"

	"github.com/aserto-dev/go-aserto/client"
	model "github.com/aserto-dev/go-directory/aserto/directory/model/v3"
	ds3 "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	dw3 "github.com/aserto-dev/go-directory/aserto/directory/writer/v3"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Client provides access to the Aserto control plane services.
type Client struct {
	conn *client.Connection

	// Reader client for directory
	Reader ds3.ReaderClient

	// Writer client for directory
	Writer dw3.WriterClient

	// Model client for directory service
	Model model.ModelClient
}

// New creates a tenant Client with the specified connection options.
func New(ctx context.Context, opts ...client.ConnectionOption) (*Client, error) {
	conn, err := client.NewConnection(ctx, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create grpc client failed")
	}

	return &Client{
		conn:   conn,
		Reader: ds3.NewReaderClient(conn.Conn),
		Writer: dw3.NewWriterClient(conn.Conn),
		Model:  model.NewModelClient(conn.Conn),
	}, err
}

// SetTenantID provides a tenantID to be included in outgoing messages.
func (c *Client) SetTenantID(tenantID string) {
	c.conn.TenantID = tenantID
}

// Connection returns the underlying grpc connection.
func (c *Client) Connection() grpc.ClientConnInterface {
	return c.conn.Conn
}
