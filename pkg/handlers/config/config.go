package config

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/go-aserto/client"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	account "github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"

	"github.com/pkg/errors"
)

type GetEnvCmd struct{}

func (cmd *GetEnvCmd) Run(c *cc.CommonCtx) error {
	return jsonx.OutputJSON(c.UI.Output(), c.Environment)
}

type GetTenantCmd struct{}

func (cmd *GetTenantCmd) Run(c *cc.CommonCtx) error {
	client, err := c.TenantClient()
	if err != nil {
		return err
	}

	req := &account.GetAccountRequest{}

	resp, err := client.Account.GetAccount(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	type tenant struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Current bool   `json:"current"`
		Default bool   `json:"default"`
	}

	tenants := make([]*tenant, len(resp.Result.Tenants))

	for i, t := range resp.Result.Tenants {
		isCurrent := (t.Id == c.TenantID())
		isDefault := (t.Id == resp.Result.DefaultTenant)
		tt := tenant{
			ID:      t.Id,
			Name:    t.Name,
			Current: isCurrent,
			Default: isDefault,
		}
		tenants[i] = &tt
	}

	return jsonx.OutputJSON(c.UI.Output(), tenants)
}

type SetTenantCmd struct {
	ID      string `arg:"" required:"" name:"tenant-id" help:"tenant id"`
	Default bool   `name:"default" help:"set default tenant for user"`
}

func (cmd *SetTenantCmd) Run(c *cc.CommonCtx) error {
	conn, err := c.TenantClient()
	if err != nil {
		return err
	}

	getAccntResp, err := conn.Account.GetAccount(c.Context, &account.GetAccountRequest{})
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	var tnt *api.Tenant
	for _, t := range getAccntResp.Result.Tenants {
		if t.Id == cmd.ID {
			tnt = t
			break
		}
	}

	if tnt == nil {
		return errors.Errorf("tenant id does not exist in users tenant collection [%s]", cmd.ID)
	}

	fmt.Fprintf(c.UI.Err(), "tenant %s - %s\n", tnt.Id, tnt.Name)

	tok, err := c.Token()
	if err != nil {
		return err
	}

	tok.TenantID = tnt.Id

	if err = user.GetConnectionKeys(c.Context, conn, tok); err != nil {
		return errors.Wrapf(err, "get connection keys")
	}

	kr, err := keyring.NewKeyRing(c.Auth.Issuer)
	if err != nil {
		return err
	}
	if err := kr.SetToken(tok); err != nil {
		return err
	}

	if cmd.Default {
		fmt.Fprintf(c.UI.Err(), "set default tenant to [%s]\n", cmd.ID)

		req := &account.UpdateAccountRequest{
			Account: &api.Account{DefaultTenant: cmd.ID},
		}

		if _, err := conn.Account.UpdateAccount(client.SetTenantContext(c.Context, tok.TenantID), req); err != nil {
			return errors.Wrapf(err, "update account")
		}
	}

	return nil
}
