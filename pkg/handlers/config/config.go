package config

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	errs "github.com/aserto-dev/aserto/pkg/cc/errors"
	"github.com/aserto-dev/aserto/pkg/handlers/user"
	"github.com/aserto-dev/go-grpc/aserto/api/v1"
	account "github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"
	topazConfig "github.com/aserto-dev/topaz/pkg/cc/config"
	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
	topazConfigure "github.com/aserto-dev/topaz/pkg/cli/cmd/configure"
	"github.com/aserto-dev/topaz/pkg/cli/table"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type ListConfigCmd struct {
	topazConfigure.ListConfigCmd
}

type tenant struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Current bool   `json:"current"`
	Default bool   `json:"default"`
}

func (cmd *ListConfigCmd) Run(c *cc.CommonCtx) error {
	tab := table.New(c.StdErr()).WithColumns("", "Name", "Config")

	resp, err := getAccountDetails(c)
	if err != nil && !errors.Is(err, errs.ErrNeedLogin) {
		return err
	}
	if resp != nil {
		tenants := make([]*tenant, len(resp.Result.Tenants))

		slices.SortFunc(resp.Result.Tenants, func(a, b *api.Tenant) int {
			return cmp.Compare(a.Name, b.Name)
		})

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
			name := fmt.Sprintf("%s%s", t.Name, cc.TenantSuffix)

			active := ""
			if c.Config.ConfigName == name {
				active = "*"
			}

			tab.WithRow(active, name, t.Id)
		}
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

		tab.WithRow(active, name, files[i].Name())
	}

	tab.Do()

	return nil
}

type UseConfigCmd struct {
	topazConfigure.UseConfigCmd
}

func (cmd *UseConfigCmd) Run(c *cc.CommonCtx) error {
	c.Config.ConfigName = string(cmd.Name)

	if !cc.IsAsertoAccount(c.Config.ConfigName) {
		c.Config.TenantID = ""

		topazUse := topazConfigure.UseConfigCmd{
			Name:      cmd.Name,
			ConfigDir: topazCC.GetTopazCfgDir(),
		}

		err := topazUse.Run(c.CommonCtx)
		if err != nil {
			return err
		}

		loader, err := topazConfig.LoadConfiguration(c.CommonCtx.Config.Active.ConfigFile)
		if err != nil {
			return err
		}
		servicesConfig := loader.Configuration.OPA.Config.Services

		serviceMap, ok := servicesConfig["aserto-discovery"].(map[string]interface{})
		if ok {
			headersMap, ok := serviceMap["headers"].(map[string]interface{})
			if ok {
				strTenantID, ok := headersMap["aserto-tenant-id"].(string)
				c.Config.TenantID = strTenantID
				if !ok {
					c.Config.TenantID = ""
				}
			}
		}

	} else {
		tenantName := strings.TrimSuffix(c.Config.ConfigName, cc.TenantSuffix)

		resp, err := getAccountDetails(c)
		if err != nil {
			return err
		}

		tenant := lo.Filter(resp.Result.Tenants, func(item *api.Tenant, index int) bool {
			return item.Name == tenantName
		})

		if len(tenant) != 1 {
			return errors.Wrapf(errs.ErrResolveTenant, tenantName) //nolint: govet
		}

		token, err := c.Token()
		if err != nil {
			return err
		}

		c.Config.TenantID = tenant[0].Id
		token.TenantID = tenant[0].Id

		if err := user.SwitchKeyRing(c, token, tenant[0].Id); err != nil {
			return err
		}
	}

	return c.SaveContextConfig(config.DefaultConfigFilePath)
}

func getAccountDetails(c *cc.CommonCtx) (*account.GetAccountResponse, error) {
	client, err := c.TenantClient(c.Context)
	if err != nil {
		return nil, err
	}

	req := &account.GetAccountRequest{}

	resp, err := client.Account.GetAccount(c.Context, req)
	if err != nil {
		return nil, errors.Wrapf(err, "get account")
	}

	return resp, nil
}
