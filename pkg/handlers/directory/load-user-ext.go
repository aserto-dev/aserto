package directory

import (
	"fmt"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dirx"
	"github.com/aserto-dev/aserto/pkg/dirx/auth0"
	jsonproducer "github.com/aserto-dev/aserto/pkg/dirx/json"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/authorizer"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/pkg/errors"
)

const (
	providerJSON  string = "json"
	providerAuth0 string = "auth0"
)

type LoadUserExtCmd struct {
	Provider string `required:"" help:"load users provider (json | auth0)" enum:"json,auth0"`
	Profile  string `optional:"" type:"existingfile" help:"provider profile file (.env)"`
	File     string `optional:"" type:"existingfile" help:"input file (.json)"`
}

func (cmd *LoadUserExtCmd) Run(c *cc.CommonCtx) error {
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
	r := make(chan *dirx.Result, 1)

	errc := make(chan error, 1)
	go func() {
		for e := range errc {
			fmt.Fprintf(c.ErrWriter, "%s\n", e.Error())
		}
	}()

	go dirx.UserExtSubscriber(ctx, dirClient, s, r, errc)

	switch cmd.Provider {
	case providerJSON:
		p := jsonproducer.NewProducer(cmd.File)
		p.Producer(s, errc)
		fmt.Fprintf(c.ErrWriter, "produced %d\n", p.Count())

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
		fmt.Fprintf(c.ErrWriter, "produced %d\n", p.Count())

	default:
		return errors.Errorf("unknown load user provider %s", cmd.Provider)
	}

	// close subscriber channel to indicate that the producer done
	close(s)

	// wait for done from subscriber, indicating last received messages has been send
	result := <-r

	// close error channel as the last action before returning
	close(errc)

	if result.Err != nil {
		return result.Err
	}

	fmt.Fprintf(c.ErrWriter, "received %d\n", result.Counts.Received)
	fmt.Fprintf(c.ErrWriter, "created  %d\n", result.Counts.Created)
	fmt.Fprintf(c.ErrWriter, "updated  %d\n", result.Counts.Updated)
	fmt.Fprintf(c.ErrWriter, "deleted  %d\n", result.Counts.Deleted)
	fmt.Fprintf(c.ErrWriter, "errors   %d\n", result.Counts.Errors)

	return nil
}
