package config

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/go-lib/ids"
	"github.com/aserto-dev/proto/aserto/api"
	"github.com/aserto-dev/proto/aserto/tenant/account"

	"github.com/pkg/errors"
)

type GetEnvCmd struct {
}

func (cmd *GetEnvCmd) Run(c *cc.CommonCtx) error {
	services, err := grpcc.Environment(c.Environment())
	if err != nil {
		return err
	}

	return jsonx.OutputJSON(c.OutWriter, services)
}

type GetTenantCmd struct {
}

func (cmd *GetTenantCmd) Run(c *cc.CommonCtx) error {

	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	accntClient := conn.AccountClient()

	req := &account.GetAccountRequest{}

	resp, err := accntClient.GetAccount(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	type tenant struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Default bool   `json:"default"`
	}

	tenants := make([]*tenant, len(resp.Result.Tenants))

	for i, t := range resp.Result.Tenants {
		isDefault := (t.Id == resp.Result.DefaultTenant)
		tt := tenant{
			ID:      t.Id,
			Name:    t.Name,
			Default: isDefault,
		}
		tenants[i] = &tt
	}

	return jsonx.OutputJSON(c.OutWriter, tenants)
}

type SetTenantCmd struct {
	ID      string `arg:"" required:"" name:"tenant-id" help:"tenant id"`
	Default bool   `name:"default" help:"set default tenant for user"`
}

func (cmd *SetTenantCmd) Run(c *cc.CommonCtx) error {
	if err := ids.CheckTenantID(cmd.ID); err != nil {
		return errors.Errorf("argument is not a valid tenant id")
	}

	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	accntClient := conn.AccountClient()

	getAccntResp, err := accntClient.GetAccount(c.Context, &account.GetAccountRequest{})
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

	fmt.Fprintf(c.ErrWriter, "tenant %s - %s\n", tnt.Id, tnt.Name)

	tok := c.Token()
	tok.TenantID = tnt.Id

	if err := user.GetConnectionKeys(c.Context, conn, tok); err != nil {
		return errors.Wrapf(err, "get connection keys")
	}

	kr, err := keyring.NewKeyRing()
	if err != nil {
		return err
	}
	if err := kr.SetToken(c.Environment(), tok); err != nil {
		return err
	}

	if cmd.Default {
		fmt.Fprintf(c.ErrWriter, "set default tenant to [%s]\n", cmd.ID)

		req := &account.UpdateAccountRequest{
			Account: &api.Account{DefaultTenant: cmd.ID},
		}

		if _, err := accntClient.UpdateAccount(c.Context, req); err != nil {
			return errors.Wrapf(err, "update account")
		}
	}

	return nil
}
