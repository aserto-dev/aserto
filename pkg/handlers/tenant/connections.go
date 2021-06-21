package tenant

import (
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/proto/aserto/api"
	"github.com/aserto-dev/proto/aserto/tenant/connection"
	"github.com/pkg/errors"
)

type ListConnectionsCmd struct {
	Kind string `help:"provider kind"`
}

func (cmd ListConnectionsCmd) Run(c *cc.CommonCtx) error {
	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())

	connClient := conn.ConnectionManagerClient()

	req := &connection.ListConnectionsRequest{}

	kindStr := strings.ToUpper(cmd.Kind)

	if kind, ok := api.ProviderKind_value[kindStr]; ok {
		req.Kind = api.ProviderKind(kind)
	} else {
		req.Kind = api.ProviderKind_UNKNOWN_PROVIDER_KIND
	}

	resp, err := connClient.ListConnections(ctx, req)
	if err != nil {
		return errors.Wrapf(err, "list connections")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type GetConnectionCmd struct {
	ID string `help:"connection id"`
}

func (cmd GetConnectionCmd) Run(c *cc.CommonCtx) error {
	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())

	connClient := conn.ConnectionManagerClient()

	req := &connection.GetConnectionRequest{
		Id: cmd.ID,
	}

	resp, err := connClient.GetConnection(ctx, req)
	if err != nil {
		return errors.Wrapf(err, "get connection [%s]", cmd.ID)
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type VerifyConnectionCmd struct {
	ID string `help:"connection id"`
}

func (cmd VerifyConnectionCmd) Run(c *cc.CommonCtx) error {
	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())

	connClient := conn.ConnectionManagerClient()

	req := &connection.VerifyConnectionRequest{
		Id: cmd.ID,
	}

	resp, err := connClient.VerifyConnection(ctx, req)
	if err != nil {
		return errors.Wrapf(err, "verify connection [%s]", cmd.ID)
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
