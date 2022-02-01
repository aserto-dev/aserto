package tenant

import (
	"strings"

	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/client/tenant"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	provider "github.com/aserto-dev/go-grpc/aserto/tenant/provider/v1"

	"github.com/pkg/errors"
)

type ListProviderKindsCmd struct{}

func (cmd ListProviderKindsCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(
		c.Context,
		aserto.WithAddr(c.TenantService()),
		aserto.WithTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	req := &provider.ListProviderKindsRequest{}

	resp, err := client.Provider.ListProviderKinds(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "list provider kinds")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type ListProvidersCmd struct {
	Kind string `help:"provider kind"`
}

func ProviderKind(kind string) api.ProviderKind {
	kind = strings.ToUpper(kind)

	if apiKind, ok := api.ProviderKind_value[kind]; ok {
		return api.ProviderKind(apiKind)
	}
	return api.ProviderKind_PROVIDER_KIND_UNKNOWN
}

func (cmd ListProvidersCmd) Run(c *cc.CommonCtx) error {
	client, err := tenant.New(
		c.Context,
		aserto.WithAddr(c.TenantService()),
		aserto.WithTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	resp, err := client.Provider.ListProviders(
		c.Context,
		&provider.ListProvidersRequest{
			Kind: ProviderKind(cmd.Kind),
		})
	if err != nil {
		return errors.Wrapf(err, "list providers")
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}

type GetProviderCmd struct {
	ID string `arg:"" required:"" help:"provider id"`
}

func (cmd GetProviderCmd) Run(c *cc.CommonCtx) error {
	conn, err := tenant.New(
		c.Context,
		aserto.WithAddr(c.TenantService()),
		aserto.WithTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	req := &provider.GetProviderRequest{
		Id: cmd.ID,
	}

	resp, err := conn.Provider.GetProvider(c.Context, req)
	if err != nil {
		return errors.Wrapf(err, "get provider [%s]", cmd.ID)
	}

	return jsonx.OutputJSONPB(c.OutWriter, resp)
}
