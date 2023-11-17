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

type ClientReader struct {
	conn *client.Connection

	// Reader client for directory
	Reader ds3.ReaderClient
}

type ClientWriter struct {
	conn *client.Connection

	// Writer client for directory
	Writer dw3.WriterClient
}

type ClientModel struct {
	conn *client.Connection

	// Model client for directory service
	Model model.ModelClient
}

// New creates a directory reader Client with the specified connection options.
func NewReader(ctx context.Context, opts ...client.ConnectionOption) (*ClientReader, error) {
	conn, err := client.NewConnection(ctx, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create grpc client failed")
	}

	return &ClientReader{
		conn:   conn,
		Reader: ds3.NewReaderClient(conn.Conn),
	}, err
}

// SetTenantID provides a tenantID to be included in outgoing messages.
func (c *ClientReader) SetTenantID(tenantID string) {
	c.conn.TenantID = tenantID
}

// Connection returns the underlying grpc connection.
func (c *ClientReader) Connection() grpc.ClientConnInterface {
	return c.conn.Conn
}

// New creates a directory writer Client with the specified connection options.
func NewWriter(ctx context.Context, opts ...client.ConnectionOption) (*ClientWriter, error) {
	conn, err := client.NewConnection(ctx, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create grpc client failed")
	}

	return &ClientWriter{
		conn:   conn,
		Writer: dw3.NewWriterClient(conn.Conn),
	}, err
}

// SetTenantID provides a tenantID to be included in outgoing messages.
func (c *ClientWriter) SetTenantID(tenantID string) {
	c.conn.TenantID = tenantID
}

// Connection returns the underlying grpc connection.
func (c *ClientWriter) Connection() grpc.ClientConnInterface {
	return c.conn.Conn
}

// New creates a directory model Client with the specified connection options.
func NewModel(ctx context.Context, opts ...client.ConnectionOption) (*ClientModel, error) {
	conn, err := client.NewConnection(ctx, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "create grpc client failed")
	}

	return &ClientModel{
		conn:  conn,
		Model: model.NewModelClient(conn.Conn),
	}, err
}

// SetTenantID provides a tenantID to be included in outgoing messages.
func (c *ClientModel) SetTenantID(tenantID string) {
	c.conn.TenantID = tenantID
}

// Connection returns the underlying grpc connection.
func (c *ClientModel) Connection() grpc.ClientConnInterface {
	return c.conn.Conn
}
