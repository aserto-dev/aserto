package config

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
	"github.com/aserto-dev/aserto/pkg/jsonx"
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
	ID string `arg:"" required:"" name:"tenant-id" help:"tenant id"`
}

func (cmd *SetTenantCmd) Run(c *cc.CommonCtx) error {
	if err := ids.CheckTenantID(cmd.ID); err != nil {
		return errors.Errorf("argument is not a valid tenant id")
	}

	fmt.Fprintf(c.ErrWriter, "tenant %s\n", cmd.ID)

	fmt.Fprintf(c.ErrWriter, "set default tenant to [%s]\n", cmd.ID)

	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	accntClient := conn.AccountClient()

	req := &account.UpdateAccountRequest{
		Account: &api.Account{DefaultTenant: cmd.ID},
	}

	if _, err := accntClient.UpdateAccount(c.Context, req); err != nil {
		return errors.Wrapf(err, "update account")
	}

	return nil
}
