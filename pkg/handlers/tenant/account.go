package tenant

import (
	"github.com/aserto-dev/aserto-go/client/tenant"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	account "github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"

	"github.com/pkg/errors"
)

type GetAccountCmd struct{}

func (cmd GetAccountCmd) Run(c *cc.CommonCtx) error {
	conn, err := tenant.New(c.Context, c.TenantSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	req := &account.GetAccountRequest{}

	resp, err := conn.Account.GetAccount(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
