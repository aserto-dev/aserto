package tenant

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/proto/aserto/tenant/policy"
	"github.com/pkg/errors"
)

type ListPolicyReferencesCmd struct{}

func (cmd ListPolicyReferencesCmd) Run(c *cc.CommonCtx) error {
	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())

	policyClient := conn.PolicyClient()

	req := &policy.ListPolicyRefsRequest{}

	resp, err := policyClient.ListPolicyRefs(ctx, req)
	if err != nil {
		return errors.Wrapf(err, "list policy packages")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
