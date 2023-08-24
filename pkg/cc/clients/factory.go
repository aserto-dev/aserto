package clients

import (
	"context"
	"log"

	"github.com/aserto-dev/aserto/pkg/cc/config"
	token_ "github.com/aserto-dev/aserto/pkg/cc/token"
	tenant_ "github.com/aserto-dev/aserto/pkg/client/tenant"
	"github.com/aserto-dev/aserto/pkg/x"
	aserto "github.com/aserto-dev/go-aserto/client"
	"github.com/aserto-dev/go-aserto/client/authorizer"
	dl "github.com/aserto-dev/go-decision-logs/aserto/decision-logs/v2"
	"github.com/aserto-dev/go-grpc/aserto/management/v2"
	"github.com/pkg/errors"
)

type Factory interface {
	TenantID() string

	TenantClient() (*tenant_.Client, error)
	AuthorizerClient() (*authorizer.Client, error)
	DecisionLogsClient() (dl.DecisionLogsClient, error)
	ControlPlaneClient() (management.ControlPlaneClient, error)
}

type OptionsBuilder func() ([]aserto.ConnectionOption, error)

type AsertoFactory struct {
	ctx        context.Context
	tenantID   string
	svcOptions map[x.Service]OptionsBuilder
}

type TenantID string

func NewClientFactory(
	ctx context.Context,
	ctxs config.Context,
	services *x.Services,
	token *token_.CachedToken,
) (*AsertoFactory, error) {
	var tenant string
	for _, cts := range ctxs.Contexts {
		if cts.Name == ctxs.ActiveContext {
			tenant = cts.TenantID
		}
	}

	defaultEnv := x.DefaultEnvironment()

	options := map[x.Service]OptionsBuilder{}
	for _, svc := range x.AllServices {
		cfg := &optionsBuilder{
			service:     svc,
			options:     services.Get(svc),
			defaultAddr: defaultEnv.Get(svc).Address,
			tenantID:    tenant,
			token:       token,
		}

		options[svc] = cfg.ConnectionOptions
	}

	for _, svc := range x.AllServices {
		cfg := &optionsBuilder{
			service:     svc,
			options:     services.Get(svc),
			defaultAddr: defaultEnv.Get(svc).Address,
			tenantID:    tenant,
			token:       token,
		}

		options[svc] = cfg.ConnectionOptions
	}

	return &AsertoFactory{
		ctx:        ctx,
		tenantID:   tenant,
		svcOptions: options,
	}, nil
}

func (c *AsertoFactory) TenantID() string {
	return c.tenantID
}

func (c *AsertoFactory) TenantClient() (*tenant_.Client, error) {
	options, err := c.options(x.TenantService)
	if err != nil {
		return nil, err
	}
	return tenant_.New(c.ctx, options...)
}

func (c *AsertoFactory) AuthorizerClient() (*authorizer.Client, error) {
	options, err := c.options(x.AuthorizerService)
	if err != nil {
		return nil, err
	}
	return authorizer.New(c.ctx, options...)
}

func (c *AsertoFactory) DecisionLogsClient() (dl.DecisionLogsClient, error) {
	options, err := c.options(x.DecisionLogsService)
	if err != nil {
		return nil, err
	}

	conn, err := aserto.NewConnection(c.ctx, options...)
	if err != nil {
		return nil, err
	}

	return dl.NewDecisionLogsClient(conn.Conn), nil
}

func (c *AsertoFactory) ControlPlaneClient() (management.ControlPlaneClient, error) {
	options, err := c.options(x.ControlPlaneService)
	if err != nil {
		return nil, err
	}

	conn, err := aserto.NewConnection(c.ctx, options...)
	if err != nil {
		return nil, err
	}

	return management.NewControlPlaneClient(conn.Conn), nil
}

func (c *AsertoFactory) options(svc x.Service) ([]aserto.ConnectionOption, error) {
	opts, ok := c.svcOptions[svc]
	if ok {
		return opts()
	}

	log.Panicf("missing setting for service [%s]", svc.Name())
	return nil, errors.Errorf("missing setting for service [%s]", svc.Name())
}
