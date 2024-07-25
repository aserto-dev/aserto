package controlplane

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/go-grpc/aserto/api/v2"
	"github.com/aserto-dev/go-grpc/aserto/management/v2"
)

type EdgeDirSyncCmd struct {
	Instance string `arg:"" help:"target instance"`
}

func (cmd EdgeDirSyncCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.ControlPlaneClient(c.Context)
	if err != nil {
		return err
	}

	_, err = cli.ExecCommand(c.Context, &management.ExecCommandRequest{
		Id: cmd.Instance,
		Command: &api.Command{
			Data: &api.Command_SyncEdgeDirectory{},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
