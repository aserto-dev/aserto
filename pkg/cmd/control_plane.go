package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	controlplane "github.com/aserto-dev/aserto/pkg/handlers/control_plane"
)

type ControlPlaneCmd struct {
	ListConnections           controlplane.ListConnectionsCmd           `cmd:"" help:"list satellite connections" group:"control-plane"`
	ClientCert                controlplane.ClientCertCmd                `cmd:"" help:"get client certificates for a satellite connection" group:"control-plane"`
	ListInstanceRegistrations controlplane.ListInstanceRegistrationsCmd `cmd:"" help:"list instance registrations" group:"control-plane"`
	Discovery                 controlplane.DiscoveryCmd                 `cmd:"" help:"run discovery on a registered instance" group:"control-plane"`
	EdgeDirSync               controlplane.EdgeDirSyncCmd               `cmd:"" help:"sync the directory on an edge authorizer" group:"control-plane"`
}

func (cmd *ControlPlaneCmd) Run(c *cc.CommonCtx) error {
	return nil
}
