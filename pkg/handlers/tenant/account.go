package tenant

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	account "github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"
	"github.com/aserto-dev/topaz/pkg/cli/jsonx"

	"github.com/pkg/errors"
)

type GetAccountCmd struct{}

func (cmd GetAccountCmd) Run(c *cc.CommonCtx) error {
	conn, err := c.TenantClient(c.Context)
	if err != nil {
		return err
	}

	req := &account.GetAccountRequest{}

	resp, err := conn.Account.GetAccount(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	return jsonx.OutputJSONPB(c.StdOut(), resp)
}
