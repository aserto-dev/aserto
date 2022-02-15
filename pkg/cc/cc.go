package cc

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aserto-dev/aserto-go/client/authorizer"
	"github.com/aserto-dev/aserto-go/client/tenant"
	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/client"
	"github.com/aserto-dev/aserto/pkg/x"
	dl "github.com/aserto-dev/go-grpc/aserto/decision_logs/v1"
)

type CommonCtx struct {
	Context       context.Context
	OutWriter     io.Writer
	ErrWriter     io.Writer
	Services      *x.Services
	clientFactory client.Factory
	tenantID      string

	token *api.Token
}

func NewCommonCtx(env *x.Services, tenantID string, clientFactory client.Factory, token *api.Token) *CommonCtx {
	return &CommonCtx{
		Context:       context.Background(),
		OutWriter:     os.Stdout,
		ErrWriter:     os.Stderr,
		Services:      env,
		tenantID:      tenantID,
		clientFactory: clientFactory,
		token:         token,
	}
}

func (ctx *CommonCtx) Environment() string {
	return ctx.Services.Environment
}

func (ctx *CommonCtx) AccessToken() string {
	return ctx.token.Access
}

func (ctx *CommonCtx) Token() *api.Token {
	return ctx.token
}

func (ctx *CommonCtx) TenantID() string {
	return ctx.tenantID
}

func (ctx *CommonCtx) TenantClient() (*tenant.Client, error) {
	return ctx.clientFactory.TenantClient(ctx.Context)
}

func (ctx *CommonCtx) AuthorizerClient() (*authorizer.Client, error) {
	return ctx.clientFactory.AuthorizerClient(ctx.Context)
}

func (ctx *CommonCtx) DecisionLogsClient() (dl.DecisionLogsClient, error) {
	return ctx.clientFactory.DecisionLogsClient(ctx.Context)
}

func (ctx *CommonCtx) AuthorizerAPIKey() string {
	return ctx.token.AuthorizerAPIKey
}

func (ctx *CommonCtx) DecisionLogsKey() string {
	return ctx.token.DecisionLogsKey
}

func (ctx *CommonCtx) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
