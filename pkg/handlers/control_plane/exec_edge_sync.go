package controlplane

import (
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/go-grpc/aserto/api/v2"
	"github.com/aserto-dev/go-grpc/aserto/management/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EdgeDirSyncCmd struct {
	Instance string `arg:"" optional:"" help:"target instance"`
	All      bool   `flag:"" help:"fan-out command to all online instances"`
	Mode     string `flag:"" default:"diff" enum:"full,diff,watermark,manifest" help:"sync mode"`
}

func (cmd EdgeDirSyncCmd) Run(c *cc.CommonCtx) error {
	cli, err := c.ControlPlaneClient(c.Context)
	if err != nil {
		return err
	}

	instances := []string{}

	cmd.Mode = strings.ToLower(cmd.Mode)

	var mode api.SyncMode

	switch cmd.Mode {
	case "full":
		mode = api.SyncMode_SYNC_MODE_FULL
	case "diff":
		mode = api.SyncMode_SYNC_MODE_DIFF
	case "watermark":
		mode = api.SyncMode_SYNC_MODE_WATERMARK
	case "manifest":
		mode = api.SyncMode_SYNC_MODE_MANIFEST
	default:
		mode = api.SyncMode_SYNC_MODE_FULL
	}

	switch {
	case cmd.All:
		resp, err := cli.ListInstanceRegistrations(c.Context, &management.ListInstanceRegistrationsRequest{})
		if err != nil {
			return err
		}

		for _, inst := range resp.Result {
			instances = append(instances, inst.Id)
		}

	case cmd.Instance != "" && !cmd.All:
		instances = append(instances, cmd.Instance)

	case cmd.Instance == "" && !cmd.All:
		return status.Errorf(codes.InvalidArgument, "no instance provided")
	}

	for _, instance := range instances {
		c.Con().Info().Msg("send edge-sync %q to %q", cmd.Mode, instance)

		_, err = cli.ExecCommand(c.Context, &management.ExecCommandRequest{
			Id: instance,
			Command: &api.Command{
				Data: &api.Command_SyncEdgeDirectory{
					SyncEdgeDirectory: &api.SyncEdgeDirectoryCommand{
						Mode: mode,
					},
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
