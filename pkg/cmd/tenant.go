package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/aserto/pkg/handlers/tenant"
)

type TenantCmd struct {
	Get  TenantGetCmd  `cmd:"" help:""`
	List TenantListCmd `cmd:"" help:""`
	Exec TenantExecCmd `cmd:"" help:""`
}

func (cmd *TenantCmd) BeforeApply(context *kong.Context) error {
	cfg, err := getConfig(context)
	if err != nil {
		return err
	}
	if !cc.IsAsertoAccount(cfg.ConfigName) && cfg.TenantID == "" {
		return errors.ErrTenantCmd
	}
	return nil
}

func (cmd *TenantCmd) Run(c *cc.CommonCtx) error {
	return nil
}

type TenantGetCmd struct {
	Account    tenant.GetAccountCmd    `cmd:"" help:"get account info"`
	Connection tenant.GetConnectionCmd `cmd:"" help:"get connection instance info"`
	Provider   tenant.GetProviderCmd   `cmd:"" help:"get provider info"`
}

type TenantListCmd struct {
	Connections      tenant.ListConnectionsCmd      `cmd:"" help:"list connections"`
	PolicyReferences tenant.ListPolicyReferencesCmd `cmd:"" help:"list policy references"`
	ProviderKinds    tenant.ListProviderKindsCmd    `cmd:"" help:"list provider kinds"`
	Providers        tenant.ListProvidersCmd        `cmd:"" help:"list providers"`
}

type TenantExecCmd struct {
	Update tenant.UpdateConnectionCmd `cmd:"" help:"update connection configuration fields"`
	Verify tenant.VerifyConnectionCmd `cmd:"" help:"verify connection settings"`
	Sync   tenant.SyncConnectionCmd   `cmd:"" help:"trigger sync of IDP connection"`
}
