package controlplane

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
	"github.com/pkg/errors"
)

type ClientCertCmd struct {
	ID string `arg:"" help:"satellite connection ID"`
}

func (cmd ClientCertCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.TenantClient()
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

	if conn.Kind != api.ProviderKind_PROVIDER_KIND_SATELLITE {
		return errors.New("not a satellite connection")
	}

	certs := conn.Config.Fields["api_cert"].GetListValue().GetValues()
	if len(certs) == 0 {
		return errors.New("invalid connection configuration")
	}

	err = jsonx.OutputJSONPB(c.UI.Output(), certs[len(certs)-1])
	if err != nil {
		return err
	}

	return nil
}
