package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/errors"
	cp "github.com/aserto-dev/aserto/pkg/handlers/control_plane"
)

type ControlPlaneCmd struct {
	Get  ControlPlaneGetCmd  `cmd:"" help:"certificates"`
	List ControlPlaneListCmd `cmd:"" help:"list connections | instances"`
	Exec ControlPlaneExecCmd `cmd:"" help:"exec discovery | edge-sync"`
}

func (cmd *ControlPlaneCmd) BeforeApply(context *kong.Context) error {
	cfg, err := getConfig(context)
	if err != nil {
		return err
	}

	if !cc.IsAsertoAccount(cfg.ConfigName) && cfg.TenantID == "" {
		return errors.ErrControlPlaneCmd
	}
	return nil
}

func (cmd *ControlPlaneCmd) Run(c *cc.CommonCtx) error {
	return nil
}

type ControlPlaneGetCmd struct {
	Certificates cp.ClientCertCmd `cmd:"" help:"get client certificates for an edge authorizer connection"`
}

type ControlPlaneListCmd struct {
	Connections cp.ListConnectionsCmd           `cmd:"" help:"list edge authorizer connections"`
	Instances   cp.ListInstanceRegistrationsCmd `cmd:"" help:"list instance registrations"`
}

type ControlPlaneExecCmd struct {
	Discovery cp.DiscoveryCmd   `cmd:"" help:"run discovery on a registered instance"`
	EdgeSync  cp.EdgeDirSyncCmd `cmd:"" help:"sync the directory on an edge authorizer"`
}
