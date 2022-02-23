package directory

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dirx"
	"github.com/aserto-dev/aserto/pkg/dirx/auth0"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	connection "github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"

	"github.com/pkg/errors"
)

// TODO : make using IDP connection explicit.
type LoadUsersCmd struct {
	Provider    string `required:"" help:"load users provider (json | auth0)" enum:"json,auth0"`
	Profile     string `optional:"" type:"existingfile" help:"provider profile file (.env)"`
	File        string `optional:"" type:"existingfile" help:"input file (.json)"`
	InclUserExt bool   `optional:"" help:"include user extensions (attributes & applications) in the base user object"`
}

func (cmd *LoadUsersCmd) Run(c *cc.CommonCtx) error {
	loader := UserLoader{Provider: cmd.Provider, Profile: cmd.Profile, File: cmd.File}
	return loader.Load(c, dirx.NewLoadUsersRequestFactory(cmd.InclUserExt))
}

func auth0ConfigFromConnection(c *cc.CommonCtx) (*auth0.Config, error) {
	cfg := auth0.Config{}

	client, err := c.TenantClient()
	if err != nil {
		return nil, err
	}

	listResp, err := client.Connections.ListConnections(
		c.Context,
		&connection.ListConnectionsRequest{
			Kind: api.ProviderKind_PROVIDER_KIND_IDP,
		},
	)
	if err != nil {
		return nil, err
	}
	if len(listResp.Results) != 1 {
		return nil, errors.Errorf("identity provider connection not found")
	}

	connResp, err := client.Connections.GetConnection(
		c.Context,
		&connection.GetConnectionRequest{
			Id: listResp.Results[0].Id,
		},
	)
	if err != nil {
		return nil, err
	}

	const (
		confKeyAuth0Domain       string = "domain"
		confKeyAuth0ClientID     string = "client_id"
		confKeyAuth0ClientSecret string = "client_secret"
	)

	var ok bool
	connConfigMap := connResp.Result.Config.AsMap()
	cfg.Domain, ok = connConfigMap[confKeyAuth0Domain].(string)
	if !ok {
		return nil, errors.Errorf("config key not found [%s]", confKeyAuth0Domain)
	}

	cfg.ClientID, ok = connConfigMap[confKeyAuth0ClientID].(string)
	if !ok {
		return nil, errors.Errorf("config key not found [%s]", confKeyAuth0ClientID)
	}

	cfg.ClientSecret, ok = connConfigMap[confKeyAuth0ClientSecret].(string)
	if !ok {
		return nil, errors.Errorf("config key not found [%s]", confKeyAuth0ClientSecret)
	}

	return &cfg, nil
}
