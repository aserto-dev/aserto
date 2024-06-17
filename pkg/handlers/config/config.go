package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	account "github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"

	"github.com/pkg/errors"

	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazConfig "github.com/aserto-dev/topaz/pkg/cli/cmd/configure"
)

type ListConfigCmd struct {
	topazConfig.ListConfigCmd
}

func (cmd *ListConfigCmd) Run(c *cc.CommonCtx) error {
	table := c.UI.Normal().WithTable("", "Name", "Config File", "Target")

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
		table.WithTableRow(active, name, "", "remote")
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

		table.WithTableRow(active, name, files[i].Name(), "local")
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
			Name:      cmd.Name,
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
