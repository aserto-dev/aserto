package config

// import (
// 	"fmt"

// 	"github.com/aserto-dev/aserto/pkg/app"
// 	"github.com/aserto-dev/aserto/pkg/flags"
// 	"github.com/aserto-dev/aserto/pkg/grpcc"
// 	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
// 	"github.com/aserto-dev/aserto/pkg/jsonx"
// 	"github.com/aserto-dev/aserto/pkg/x"
// 	"github.com/aserto-dev/go-lib/ids"
// 	"github.com/aserto-dev/proto/aserto/api"
// 	"github.com/aserto-dev/proto/aserto/tenant/account"
// 	"github.com/pkg/errors"
// 	"github.com/urfave/cli/v2"
// )

// func GetEnvHandler(c *cli.Context) error {
// 	services := grpcc.Environment(c.String(flags.FlagEnvironment))
// 	fmt.Fprintf(c.App.ErrWriter, "%s - default env [%s]\n", x.Target, x.DefaultEnv)
// 	return jsonx.OutputJSON(c.App.Writer, services)
// }

// func SetEnvHandler(c *cli.Context) error {
// 	fmt.Fprintf(c.App.ErrWriter, "set env handler\n")
// 	return nil
// }

// func GetTenantHandler(c *cli.Context) error {
// 	appCtx := app.GetAppContext(c.Context)

// 	conn, err := tenant.Connection(
// 		c.Context,
// 		appCtx.TenantService(),
// 		grpcc.NewTokenAuth(appCtx.AccessToken()),
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	accntClient := conn.AccountClient()

// 	req := &account.GetAccountRequest{}

// 	resp, err := accntClient.GetAccount(c.Context, req)
// 	if err != nil {
// 		return errors.Wrapf(err, "get account")
// 	}

// 	type tenant struct {
// 		ID      string `json:"id"`
// 		Name    string `json:"name"`
// 		Default bool   `json:"default"`
// 	}

// 	tenants := make([]*tenant, len(resp.Result.Tenants))

// 	for i, t := range resp.Result.Tenants {
// 		isDefault := (t.Id == resp.Result.DefaultTenant)
// 		tt := tenant{
// 			ID:      t.Id,
// 			Name:    t.Name,
// 			Default: isDefault,
// 		}
// 		tenants[i] = &tt
// 	}

// 	return jsonx.OutputJSON(c.App.Writer, tenants)
// }

// func SetTenantHandler(c *cli.Context) error {
// 	appCtx := app.GetAppContext(c.Context)

// 	if !c.Args().Present() || !(c.Args().Len() <= 2) {
// 		return errors.Errorf("invalid number of arguments")
// 	}

// 	id := c.Args().First()
// 	if err := ids.CheckTenantID(id); err != nil {
// 		return errors.Errorf("argument is not a valid tenant id")
// 	}

// 	fmt.Fprintf(c.App.ErrWriter, "tenant %s\n", id)

// 	setDefault := flags.GetBoolTailFlag(c, "default")
// 	if setDefault {
// 		fmt.Fprintf(c.App.ErrWriter, "set default tenant to [%s]\n", id)

// 		conn, err := tenant.Connection(
// 			c.Context,
// 			appCtx.TenantService(),
// 			grpcc.NewTokenAuth(appCtx.AccessToken()),
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		accntClient := conn.AccountClient()

// 		req := &account.UpdateAccountRequest{
// 			Account: &api.Account{DefaultTenant: id},
// 		}

// 		if _, err := accntClient.UpdateAccount(c.Context, req); err != nil {
// 			return errors.Wrapf(err, "update account")
// 		}
// 	}

// 	return nil
// }
