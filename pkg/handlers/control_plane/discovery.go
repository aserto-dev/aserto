package controlplane

import (
	"github.com/aserto-dev/aserto/pkg/api/management/v2"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/go-grpc/aserto/api/v2"
)

type DiscoveryCmd struct {
	Instance string `arg:"" help:"target instance"`
}

func (cmd DiscoveryCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.ControlPlaneClient()
	if err != nil {
		return err
	}

	_, err = cli.ExecCommand(c.Context, &management.ExecCommandRequest{
		Id: cmd.Instance,
		Command: &api.Command{
			Data: &api.Command_Discovery{},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
