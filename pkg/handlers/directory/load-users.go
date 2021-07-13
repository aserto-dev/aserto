package directory

import (
	"context"
	"fmt"

	"github.com/aserto-dev/aserto-tenant/pkg/app/providers"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dirx"
	"github.com/aserto-dev/aserto/pkg/dirx/auth0"
	jsonproducer "github.com/aserto-dev/aserto/pkg/dirx/json"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
	"github.com/aserto-dev/proto/aserto/api"
	"github.com/aserto-dev/proto/aserto/tenant/connection"

	"github.com/pkg/errors"
)

// TODO : make using IDP connection explicit
type LoadUsersCmd struct {
	Provider    string `required:"" help:"load users provider (json | auth0)" enum:"json,auth0"`
	Profile     string `optional:"" type:"existingfile" help:"provider profile file (.env)"`
	File        string `optional:"" type:"existingfile" help:"input file (.json)"`
	InclUserExt bool   `optional:"" help:"include user extensions (attributes & applications) in the base user object"`
}

func (cmd *LoadUsersCmd) Run(c *cc.CommonCtx) error {
	conn, err := authorizer.Connection(
		c.Context,
		c.AuthorizerService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return err
	}

	ctx := grpcc.SetTenantContext(c.Context, c.TenantID())
	dirClient := conn.DirectoryClient()

	s := make(chan *api.User, 10)
	done := make(chan bool, 1)

	errc := make(chan error, 1)
	go func() {
		for e := range errc {
			fmt.Fprintf(c.ErrWriter, "%s\n", e.Error())
		}
	}()

	go dirx.Subscriber(ctx, dirClient, s, done, errc, cmd.InclUserExt)

	switch cmd.Provider {
	case providerJSON:
		p := jsonproducer.NewProducer(cmd.File)
		p.Producer(s, errc)
		fmt.Fprintf(c.ErrWriter, "produced %d instances\n", p.Count())

	case providerAuth0:
		var cfg *auth0.Config

		if cmd.Profile != "" {
			cfg, err = auth0.FromProfile(cmd.Profile)
			if err != nil {
				return err
			}
		} else {
			cfg, err = auth0ConfigFromConnection(c)
			if err != nil {
				return err
			}
		}

		if err := cfg.Validate(); err != nil {
			return err
		}

		p := auth0.NewProducer(cfg)
		p.Producer(s, errc)
		fmt.Fprintf(c.ErrWriter, "produced %d instances\n", p.Count())

	default:
		return errors.Errorf("unknown load user provider %s", cmd.Provider)
	}

	// close subscriber channel to indicate that the producer done
	close(s)

	// wait for done from subscriber, indicating last received messages has been send
	<-done

	// close error channel as the last action before returning
	close(errc)

	return nil
}

func auth0ConfigFromConnection(c *cc.CommonCtx) (*auth0.Config, error) {
	cfg := auth0.Config{}

	conn, err := tenant.Connection(
		c.Context,
		c.TenantService(),
		grpcc.NewTokenAuth(c.AccessToken()),
	)
	if err != nil {
		return nil, err
	}

	ctx := grpcc.SetTenantContext(context.Background(), c.TenantID())

	connClient := conn.ConnectionManagerClient()
	listResp, err := connClient.ListConnections(
		ctx,
		&connection.ListConnectionsRequest{
			Kind: api.ProviderKind_IDP,
		},
	)
	if err != nil {
		return nil, err
	}
	if len(listResp.Results) != 1 {
		return nil, errors.Errorf("identity provider connection not found")
	}

	connResp, err := connClient.GetConnection(
		ctx,
		&connection.GetConnectionRequest{
			Id: listResp.Results[0].Id,
		},
	)
	if err != nil {
		return nil, err
	}

	var ok bool
	connConfigMap := connResp.Result.Config.AsMap()
	cfg.Domain, ok = connConfigMap[providers.ConfKeyAuth0Domain].(string)
	if !ok {
		return nil, errors.Errorf("config key not found [%s]", providers.ConfKeyAuth0Domain)
	}

	cfg.ClientID, ok = connConfigMap[providers.ConfKeyAuth0ClientID].(string)
	if !ok {
		return nil, errors.Errorf("config key not found [%s]", providers.ConfKeyAuth0ClientID)
	}

	cfg.ClientSecret, ok = connConfigMap[providers.ConfKeyAuth0ClientSecret].(string)
	if !ok {
		return nil, errors.Errorf("config key not found [%s]", providers.ConfKeyAuth0ClientSecret)
	}

	return &cfg, nil
}
