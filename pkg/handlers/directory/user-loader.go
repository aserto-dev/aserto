package directory

import (
	"fmt"

	"github.com/aserto-dev/aserto-go/client/authorizer"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/dirx"
	"github.com/aserto-dev/aserto/pkg/dirx/auth0"
	jsonproducer "github.com/aserto-dev/aserto/pkg/dirx/json"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/pkg/errors"
)

type UserLoader struct {
	Provider string
	Profile  string
	File     string
}

func (userLoader *UserLoader) Load(c *cc.CommonCtx, requestFactory dirx.LoadUsersRequestFactory) error {
	client, err := authorizer.New(c.Context, c.AuthorizerSvcConnectionOptions()...)
	if err != nil {
		return err
	}

	userSubscriber := dirx.UserSubscriber{
		Ctx:           c.Context,
		DirClient:     client.Directory,
		SourceChannel: make(chan *api.User, 10),
		ResultChannel: make(chan *dirx.Result, 1),
		ErrorChannel:  make(chan error, 1),
	}
	go func() {
		for e := range userSubscriber.ErrorChannel {
			fmt.Fprintf(c.ErrWriter, "%s\n", e.Error())
		}
	}()

	go userSubscriber.Subscribe(requestFactory)

	switch userLoader.Provider {
	case providerJSON:
		p := jsonproducer.NewProducer(userLoader.File)
		p.Producer(userSubscriber.SourceChannel, userSubscriber.ErrorChannel)
		fmt.Fprintf(c.ErrWriter, "produced %d\n", p.Count())

	case providerAuth0:
		var cfg *auth0.Config

		if userLoader.Profile != "" {
			cfg, err = auth0.FromProfile(userLoader.Profile)
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
		p.Producer(userSubscriber.SourceChannel, userSubscriber.ErrorChannel)
		fmt.Fprintf(c.ErrWriter, "produced %d\n", p.Count())

	default:
		return errors.Errorf("unknown load user provider %s", userLoader.Provider)
	}

	// close subscriber channel to indicate that the producer done
	close(userSubscriber.SourceChannel)

	// wait for done from subscriber, indicating last received messages has been send
	result := <-userSubscriber.ResultChannel

	// close error channel as the last action before returning
	close(userSubscriber.ErrorChannel)

	fmt.Fprintf(c.ErrWriter, "received %d\n", result.Counts.Received)
	fmt.Fprintf(c.ErrWriter, "created  %d\n", result.Counts.Created)
	fmt.Fprintf(c.ErrWriter, "updated  %d\n", result.Counts.Updated)
	fmt.Fprintf(c.ErrWriter, "deleted  %d\n", result.Counts.Deleted)
	fmt.Fprintf(c.ErrWriter, "errors   %d\n", result.Counts.Errors)

	return result.Err
}
