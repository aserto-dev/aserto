package controlplane

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/pb"
	"github.com/aserto-dev/go-grpc/aserto/api/v2"
	"github.com/aserto-dev/go-grpc/aserto/management/v2"

	"github.com/samber/lo"
)

type ListInstanceRegistrationsCmd struct {
	Connection string `args:"" optional:"" help:"filter on connection ID"`
}

func (cmd ListInstanceRegistrationsCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.ControlPlaneClient(c.Context)
	if err != nil {
		return err
	}

	resp, err := cli.ListInstanceRegistrations(c.Context, &management.ListInstanceRegistrationsRequest{})
	if err != nil {
		return err
	}

	results := lo.FilterMap(resp.Result, func(x *api.InstanceRegistration, _ int) (*api.InstanceRegistration, bool) {
		if cmd.Connection == "" {
			return x, true
		}
		if x.Info.ConnectionId == cmd.Connection {
			return x, true
		}
		return nil, false
	})

	return pb.WriteMsgArray(c.StdOut(), results)
}
