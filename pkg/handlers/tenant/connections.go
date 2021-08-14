package tenant

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/proto/aserto/api"
	"github.com/aserto-dev/proto/aserto/tenant/connection"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
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
	ID string `arg:"" required:"" help:"connection id"`
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
	ID string `arg:"" required:"" help:"connection id"`
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

	if _, err = connClient.VerifyConnection(ctx, req); err != nil {
		st := status.Convert(err)
		re := regexp.MustCompile(`\r?\n`)

		fmt.Fprintf(c.ErrWriter, "verification    : failed\n")
		fmt.Fprintf(c.ErrWriter, "code            : %d\n", st.Code())
		fmt.Fprintf(c.ErrWriter, "message         : %s\n",
			re.ReplaceAllString(st.Message(), " | "))
		fmt.Fprintf(c.ErrWriter, "error           : %s\n",
			re.ReplaceAllString(st.Err().Error(), " | "))

		for _, detail := range st.Details() {
			if t, ok := detail.(*errdetails.ErrorInfo); ok {
				fmt.Fprintf(c.ErrWriter, "domain          : %s\n", t.Domain)
				fmt.Fprintf(c.ErrWriter, "reason          : %s\n", t.Reason)

				for k, v := range t.Metadata {
					fmt.Fprintf(c.ErrWriter, "detail          : %s (%s)\n", v, k)
				}
			}
		}
	} else {
		fmt.Fprintf(c.ErrWriter, "verification: succeeded\n")
	}

	return nil
}

type SyncConnectionCmd struct {
	ID string `arg:"" required:"" help:"connection id"`
}

func (cmd SyncConnectionCmd) Run(c *cc.CommonCtx) error {
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

	getReq := &connection.GetConnectionRequest{
		Id: cmd.ID,
	}

	curConn, err := connClient.GetConnection(ctx, getReq)
	if err != nil {
		return errors.Wrapf(err, "get connection [%s]", cmd.ID)
	}

	if curConn.Result.Kind != api.ProviderKind_IDP {
		return errors.Errorf("connection must be of kind IDP (provided %s)", curConn.Result.Kind.Enum().String())
	}

	updReq := &connection.UpdateConnectionRequest{
		Connection: curConn.Result,
		Force:      false,
	}

	if _, err = connClient.UpdateConnection(ctx, updReq); err != nil {
		return errors.Wrapf(err, "update connection [%s]", cmd.ID)
	}

	return nil
}
