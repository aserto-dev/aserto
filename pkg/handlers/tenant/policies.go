package tenant

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	policy "github.com/aserto-dev/go-grpc/aserto/tenant/policy/v1"
	"github.com/aserto-dev/topaz/pkg/cli/jsonx"

	"github.com/pkg/errors"
)

type ListPolicyReferencesCmd struct{}

func (cmd ListPolicyReferencesCmd) Run(c *cc.CommonCtx) error {
	client, err := c.TenantClient(c.Context)
	if err != nil {
		return err
	}

	req := &policy.ListPolicyRefsRequest{}

	resp, err := client.Policy.ListPolicyRefs(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "list policy packages")
	}

	return jsonx.OutputJSONPB(c.StdOut(), resp)
}
