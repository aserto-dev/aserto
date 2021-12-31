package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/tenant"
)

type TenantCmd struct {
	GetAccount           tenant.GetAccountCmd           `cmd:"" group:"tenant" help:"get account info"`
	ListConnections      tenant.ListConnectionsCmd      `cmd:"" group:"tenant" help:"list connections"`
	GetConnection        tenant.GetConnectionCmd        `cmd:"" group:"tenant" help:"get connection instance info"`
	UpdateConnection     tenant.UpdateConnectionCmd     `cmd:"" group:"tenant" help:"update connection configuration fields"`
	VerifyConnection     tenant.VerifyConnectionCmd     `cmd:"" group:"tenant" help:"verify connection settings"`
	SyncConnection       tenant.SyncConnectionCmd       `cmd:"" group:"tenant" help:"trigger sync of IDP connection"`
	ListPolicyReferences tenant.ListPolicyReferencesCmd `cmd:"" group:"tenant" help:"list policy references"`
	CreatePolicyPushKey  tenant.CreatePolicyPushKeyCmd  `cmd:"" group:"tenant" help:"create policy upload key"`
	ListProviderKinds    tenant.ListProviderKindsCmd    `cmd:"" group:"tenant" help:"list provider kinds"`
	ListProviders        tenant.ListProvidersCmd        `cmd:"" group:"tenant" help:"list providers"`
	GetProvider          tenant.GetProviderCmd          `cmd:"" group:"tenant" help:"get provider info"`
}

func (cmd *TenantCmd) BeforeApply(c *CLI) error {
	c.RequireLogin()
	return nil
}

func (cmd *TenantCmd) Run(c *cc.CommonCtx) error {
	return nil
}
