package clients

import (
	"context"
	"log"
	"time"

	token_ "github.com/aserto-dev/aserto/pkg/cc/token"
	tenant_ "github.com/aserto-dev/aserto/pkg/client/tenant"
	"github.com/aserto-dev/aserto/pkg/x"
	aserto "github.com/aserto-dev/go-aserto/client"
	dl "github.com/aserto-dev/go-decision-logs/aserto/decision-logs/v2"
	"github.com/aserto-dev/go-grpc/aserto/management/v2"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Factory interface {
	TenantID() string

	TenantClient() (*tenant_.Client, error)
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
	services *x.Services,
	tenantID TenantID,
	token *token_.CachedToken,
) (*AsertoFactory, error) {
	tenant := string(tenantID)

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

func (c *AsertoFactory) DecisionLogsClient() (dl.DecisionLogsClient, error) {
	options, err := c.options(x.DecisionLogsService)
	if err != nil {
		return nil, err
	}

	kacp := keepalive.ClientParameters{
		Time:    30 * time.Second, // send pings every 30 seconds if there is no activity
		Timeout: 5 * time.Second,  // wait 5 seconds for ping ack before considering the connection dead
	}
	options = append(options, aserto.WithDialOptions(grpc.WithKeepaliveParams(kacp)))

	conn, err := aserto.NewConnection(c.ctx, options...)
	if err != nil {
		return nil, err
	}

	return dl.NewDecisionLogsClient(conn), nil
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

	return management.NewControlPlaneClient(conn), nil
}

func (c *AsertoFactory) options(svc x.Service) ([]aserto.ConnectionOption, error) {
	opts, ok := c.svcOptions[svc]
	if ok {
		return opts()
	}

	log.Panicf("missing setting for service [%s]", svc.Name())
	return nil, errors.Errorf("missing setting for service [%s]", svc.Name())
}
