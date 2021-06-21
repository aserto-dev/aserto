package tenant

import (
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/proto/aserto/api"
	"github.com/aserto-dev/proto/aserto/tenant/provider"
	"github.com/pkg/errors"
)

type ListProviderKindsCmd struct{}

func (cmd ListProviderKindsCmd) Run(c *cc.CommonCtx) error {
	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	provClient := conn.ProviderClient()

	req := &provider.ListProviderKindsRequest{}

	resp, err := provClient.ListProviderKinds(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "list provider kinds")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type ListProvidersCmd struct {
	Kind string `help:"provider kind"`
}

func (cmd ListProvidersCmd) Run(c *cc.CommonCtx) error {
	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	provClient := conn.ProviderClient()

	req := &provider.ListProvidersRequest{}

	kindStr := strings.ToUpper(cmd.Kind)

	if kind, ok := api.ProviderKind_value[kindStr]; ok {
		req.Kind = api.ProviderKind(kind)
	} else {
		req.Kind = api.ProviderKind_UNKNOWN_PROVIDER_KIND
	}

	resp, err := provClient.ListProviders(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "list providers")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type GetProviderCmd struct {
	ID string `help:"provider id"`
}

func (cmd GetProviderCmd) Run(c *cc.CommonCtx) error {
	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	provClient := conn.ProviderClient()

	req := &provider.GetProviderRequest{
		Id: cmd.ID,
	}

	resp, err := provClient.GetProvider(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "get provider [%s]", cmd.ID)
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
