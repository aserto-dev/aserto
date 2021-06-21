package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/tenant"
)

type TenantCmd struct {
	GetAccount           tenant.GetAccountCmd           `cmd:"" group:"tenant" help:"get account info"`
	ListConnections      tenant.ListConnectionsCmd      `cmd:"" group:"tenant" help:"list connections"`
	GetConnection        tenant.GetConnectionCmd        `cmd:"" group:"tenant" help:"get connection instance info"`
	VerifyConnection     tenant.VerifyConnectionCmd     `cmd:"" group:"tenant" help:"verify connection settings"`
	ListPolicyReferences tenant.ListPolicyReferencesCmd `cmd:"" group:"tenant" help:"list policy references"`
	ListProviderKinds    tenant.ListProviderKindsCmd    `cmd:"" group:"tenant" help:"list provider kinds"`
	ListProviders        tenant.ListProvidersCmd        `cmd:"" group:"tenant" help:"list providers"`
	GetProvider          tenant.GetProviderCmd          `cmd:"" group:"tenant" help:"get provider info"`
}

func (cmd *TenantCmd) BeforeApply(c *cc.CommonCtx) error {
	return c.VerifyLoggedIn()
}

func (cmd *TenantCmd) Run(c *cc.CommonCtx) error {
	return nil
}
