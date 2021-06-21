package tenant

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/proto/aserto/tenant/account"
	"github.com/pkg/errors"
)

type GetAccountCmd struct{}

func (cmd GetAccountCmd) Run(c *cc.CommonCtx) error {
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

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
