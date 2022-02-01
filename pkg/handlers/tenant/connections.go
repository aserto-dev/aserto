package tenant

import (
	"fmt"
	"regexp"

	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/client/tenant"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	connection "github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"

	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type ListConnectionsCmd struct {
	Kind string `help:"provider kind"`
}

func (cmd ListConnectionsCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(
		c.Context,
		aserto.WithAddr(c.TenantService()),
		aserto.WithTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	resp, err := client.Connections.ListConnections(
		c.Context,
		&connection.ListConnectionsRequest{
			Kind: ProviderKind(cmd.Kind),
		})
	if err != nil {
		return errors.Wrapf(err, "list connections")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type GetConnectionCmd struct {
	ID string `arg:"" required:"" help:"connection id"`
}

func (cmd GetConnectionCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(
		c.Context,
		aserto.WithAddr(c.TenantService()),
		aserto.WithTokenAuth(c.AccessToken()),
		aserto.WithTenantID(c.TenantID()),
	)
	if err != nil {
		return err
	}

	req := &connection.GetConnectionRequest{
		Id: cmd.ID,
	}

	resp, err := client.Connections.GetConnection(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "get connection [%s]", cmd.ID)
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type VerifyConnectionCmd struct {
	ID string `arg:"" required:"" help:"connection id"`
}

func (cmd VerifyConnectionCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(
		c.Context,
		aserto.WithAddr(c.TenantService()),
		aserto.WithTokenAuth(c.AccessToken()),
		aserto.WithTenantID(c.TenantID()),
	)
	if err != nil {
		return err
	}

	req := &connection.VerifyConnectionRequest{
		Id: cmd.ID,
	}

	if _, err = client.Connections.VerifyConnection(c.Context, req); err != nil {
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
	client, err := tenant.New(
		c.Context,
		aserto.WithAddr(c.TenantService()),
		aserto.WithTokenAuth(c.AccessToken()),
		aserto.WithTenantID(c.TenantID()),
	)
	if err != nil {
		return err
	}

	getReq := &connection.GetConnectionRequest{
		Id: cmd.ID,
	}

	curConn, err := client.Connections.GetConnection(c.Context, getReq)
	if err != nil {
		return errors.Wrapf(err, "get connection [%s]", cmd.ID)
	}

	if curConn.Result.Kind != api.ProviderKind_PROVIDER_KIND_IDP {
		return errors.Errorf("connection must be of kind IDP (provided %s)", curConn.Result.Kind.Enum().String())
	}

	updReq := &connection.UpdateConnectionRequest{
		Connection: curConn.Result,
		Force:      false,
	}

	if _, err = client.Connections.UpdateConnection(c.Context, updReq); err != nil {
		return errors.Wrapf(err, "update connection [%s]", cmd.ID)
	}

	return nil
}
