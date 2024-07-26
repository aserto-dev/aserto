package controlplane

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
	"github.com/aserto-dev/topaz/pkg/cli/jsonx"

	"github.com/pkg/errors"
)

type ClientCertCmd struct {
	ID string `arg:"" help:"edge authorizer connection ID"`
}

func (cmd ClientCertCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.TenantClient(c.Context)
	if err != nil {
		return err
	}

	resp, err := cli.Connections.GetConnection(c.Context, &connection.GetConnectionRequest{
		Id: cmd.ID,
	})
	if err != nil {
		return err
	}

	conn := resp.Result
	if conn == nil {
		return errors.New("invalid empty connection")
	}

	if conn.Kind != api.ProviderKind_PROVIDER_KIND_EDGE_AUTHORIZER {
		return errors.New("not an edge authorizer connection")
	}

	certs := conn.Config.Fields["api_cert"].GetListValue().GetValues()
	if len(certs) == 0 {
		return errors.New("invalid connection configuration")
	}

	err = jsonx.OutputJSONPB(c.StdOut(), certs[len(certs)-1])
	if err != nil {
		return err
	}

	return nil
}
