package controlplane

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ListConnectionsCmd struct{}

func (cmd ListConnectionsCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.TenantClient()
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
		return errors.New("no edge authorizer connections")
	}

	var connsOut []protoreflect.ProtoMessage
	for _, conn := range conns {
		connsOut = append(connsOut, conn)
	}

	err = jsonx.OutputJSONPBArray(c.StdOut(), connsOut)
	if err != nil {
		return err
	}

	return nil
}
