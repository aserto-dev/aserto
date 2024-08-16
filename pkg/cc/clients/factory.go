package clients

import (
	"context"
	"log"
	"time"

	tok "github.com/aserto-dev/aserto/pkg/cc/token"
	cp "github.com/aserto-dev/aserto/pkg/clients/controlplane"
	dl "github.com/aserto-dev/aserto/pkg/clients/decisionlogs"
	"github.com/aserto-dev/aserto/pkg/clients/tenant"
	"github.com/aserto-dev/aserto/pkg/x"
	client "github.com/aserto-dev/go-aserto"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Factory interface {
	TenantClient(ctx context.Context) (*tenant.Client, error)
	DecisionLogsClient(ctx context.Context) (*dl.Client, error)
	ControlPlaneClient(ctx context.Context) (*cp.Client, error)
}

type OptionsBuilder func() ([]client.ConnectionOption, error)

type AsertoFactory struct {
	svcOptions map[x.Service]OptionsBuilder
}

func NewClientFactory(
	services *x.Services,
	token *tok.CachedToken,
) (*AsertoFactory, error) {
	defaultEnv := x.DefaultEnvironment()

	options := map[x.Service]OptionsBuilder{}
	for _, svc := range x.AllServices {
		cfg := &optionsBuilder{
			service:     svc,
			options:     services.Get(svc),
			defaultAddr: defaultEnv.Get(svc).Address,
			token:       token,
		}

		options[svc] = cfg.ConnectionOptions
	}

	return &AsertoFactory{
		svcOptions: options,
	}, nil
}

func (c *AsertoFactory) TenantClient(ctx context.Context) (*tenant.Client, error) {
	options, err := c.options(x.TenantService)
	if err != nil {
		return nil, err
	}

	return tenant.NewClient(ctx, options...)
}

func (c *AsertoFactory) DecisionLogsClient(ctx context.Context) (*dl.Client, error) {
	options, err := c.options(x.DecisionLogsService)
	if err != nil {
		return nil, err
	}

	kacp := keepalive.ClientParameters{
		Time:    30 * time.Second, // send pings every 30 seconds if there is no activity
		Timeout: 5 * time.Second,  // wait 5 seconds for ping ack before considering the connection dead
	}
	options = append(options, client.WithDialOptions(grpc.WithKeepaliveParams(kacp)))

	return dl.NewClient(ctx, options...)
}

func (c *AsertoFactory) ControlPlaneClient(ctx context.Context) (*cp.Client, error) {
	options, err := c.options(x.ControlPlaneService)
	if err != nil {
		return nil, err
	}

	return cp.NewClient(ctx, options...)
}

func (c *AsertoFactory) options(svc x.Service) ([]client.ConnectionOption, error) {
	opts, ok := c.svcOptions[svc]
	if ok {
		return opts()
	}

	log.Panicf("missing setting for service [%s]", svc.Name())
	return nil, errors.Errorf("missing setting for service [%s]", svc.Name())
}
