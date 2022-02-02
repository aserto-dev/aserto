package tenant

import (
	"github.com/aserto-dev/aserto-go/client/tenant"
	"github.com/aserto-dev/aserto-tenant/pkg/app/sources"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	policy "github.com/aserto-dev/go-grpc/aserto/tenant/policy/v1"

	"github.com/pkg/errors"
)

type ListPolicyReferencesCmd struct{}

func (cmd ListPolicyReferencesCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(c.Context, c.TenantSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	req := &policy.ListPolicyRefsRequest{}

	resp, err := client.Policy.ListPolicyRefs(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "list policy packages")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type CreatePolicyPushKeyCmd struct {
	PolicyID string `arg:"" required:"" help:"policy id"`
}

func (cmd CreatePolicyPushKeyCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(c.Context, c.TenantSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	req := &policy.ListPolicyRefsRequest{}
	resp, err := client.Policy.ListPolicyRefs(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "list policy packages")
	}

	found := false
	for _, p := range resp.Results {
		if found = (p.Id == cmd.PolicyID); found {
			break
		}
	}

	if !found {
		return errors.Errorf("policy id [%s] not found", cmd.PolicyID)
	}

	secret := sources.RepoPolicySecret{
		PushKey:  c.RegistryUploadKey(),
		TenantID: c.TenantID(),
		PolicyID: cmd.PolicyID,
	}

	return jsonx.OutputJSON(c.OutWriter, secret)
}
