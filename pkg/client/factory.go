package client

import (
	"context"
	"log"

	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/client/authorizer"
	"github.com/aserto-dev/aserto-go/client/tenant"
	"github.com/aserto-dev/aserto/pkg/x"
	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"
	"github.com/pkg/errors"
)

type Factory interface {
	TenantClient(ctx context.Context) (*tenant.Client, error)
	AuthorizerClient(ctx context.Context) (*authorizer.Client, error)
	DecisionLogsClient(ctx context.Context) (dl.DecisionLogsClient, error)
}

type OptionsBuilder func() ([]aserto.ConnectionOption, error)

type AsertoFactory struct {
	SvcOptions map[x.Service]OptionsBuilder
}

func (c *AsertoFactory) TenantClient(ctx context.Context) (*tenant.Client, error) {
	options, err := c.options(x.TenantService)
	if err != nil {
		return nil, err
	}
	return tenant.New(ctx, options...)
}

func (c *AsertoFactory) AuthorizerClient(ctx context.Context) (*authorizer.Client, error) {
	options, err := c.options(x.AuthorizerService)
	if err != nil {
		return nil, err
	}
	return authorizer.New(ctx, options...)
}

func (c *AsertoFactory) DecisionLogsClient(ctx context.Context) (dl.DecisionLogsClient, error) {
	options, err := c.options(x.DecisionLogsService)
	if err != nil {
		return nil, err
	}

	conn, err := aserto.NewConnection(ctx, options...)
	if err != nil {
		return nil, err
	}

	return dl.NewDecisionLogsClient(conn.Conn), nil
}

func (c *AsertoFactory) options(svc x.Service) ([]aserto.ConnectionOption, error) {
	opts, ok := c.SvcOptions[svc]
	if ok {
		return opts()
	}

	log.Panicf("missing setting for service [%s]", svc.Name())
	return nil, errors.Errorf("missing setting for service [%s]", svc.Name())
}
