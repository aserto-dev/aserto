package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/go-aserto/client"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	account "github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"

	"github.com/pkg/errors"

	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazConfig "github.com/aserto-dev/topaz/pkg/cli/cmd/configure"
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

type ListConfigCmd struct {
	topazConfig.ListConfigCmd
}

func (cmd *ListConfigCmd) Run(c *cc.CommonCtx) error {
	table := c.UI.Normal().WithTable("", "Name", "Config File")

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
		name := fmt.Sprintf("%s.%s", t.Name, t.Id)
		active := ""
		if c.Config.ConfigName == name {
			active = "*"
		}
		table.WithTableRow(active, name, "")
	}

	files, err := os.ReadDir(cmd.ConfigDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for i := range files {
		name := strings.Split(files[i].Name(), ".")[0]
		active := ""
		if name == c.Config.ConfigName {
			active = "*"
		}

		table.WithTableRow(active, name, files[i].Name())
	}
	table.Do()

	return nil
}

type UseConfigCmd struct {
	topazConfig.UseConfigCmd
}

func (cmd *UseConfigCmd) Run(c *cc.CommonCtx) error {
	c.Config.ConfigName = string(cmd.Name)
	if !cc.IsAsertoAccount(c.Config.ConfigName) {
		topazUse := topazConfig.UseConfigCmd{
			Name:      topazConfig.ConfigName(cmd.Name),
			ConfigDir: topazCC.GetTopazCfgDir(),
		}
		err := topazUse.Run(c.TopazContext)
		if err != nil {
			return err
		}
	} else {
		tenantDetails := strings.Split(c.Config.ConfigName, ".")
		c.Config.TenantID = tenantDetails[1]
	}

	return c.SaveContextConfig(config.DefaultConfigFilePath)
}
