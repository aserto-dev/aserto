package controlplane

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/pb"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
)

type ListConnectionsCmd struct{}

func (cmd ListConnectionsCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.TenantClient(c.Context)
	if err != nil {
		return err
	}

	resp, err := cli.Connections.ListConnections(c.Context, &connection.ListConnectionsRequest{
		Kind: api.ProviderKind_PROVIDER_KIND_EDGE_AUTHORIZER,
	})
	if err != nil {
		return err
	}

	conns := resp.Results
	if len(conns) == 0 {
		c.Con().Info().Msg("no edge authorizer connections")
		return nil
	}

	return pb.WriteMsgArray(c.StdOut(), conns)
}
