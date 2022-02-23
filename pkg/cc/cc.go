package cc

import (
	"context"
	"log"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/aserto-dev/clui"
)

type CommonCtx struct {
	clients.Factory

	Context     context.Context
	Environment *x.Services
	Auth        *auth0.Settings
	CachedToken token.CachedToken

	UI *clui.UI
}

func (ctx *CommonCtx) AccessToken() string {
	return ctx.Token().Access
}

func (ctx *CommonCtx) Token() *api.Token {
	return ctx.CachedToken.Get()
}

func (ctx *CommonCtx) AuthorizerAPIKey() string {
	return ctx.Token().AuthorizerAPIKey
}

func (ctx *CommonCtx) DecisionLogsKey() string {
	return ctx.Token().DecisionLogsKey
}

func (ctx *CommonCtx) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
