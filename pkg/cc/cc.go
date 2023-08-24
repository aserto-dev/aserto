package cc

import (
	"context"
	"log"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	decisionlogger "github.com/aserto-dev/aserto/pkg/decision_logger"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/aserto-dev/clui"
)

type CommonCtx struct {
	clients.Factory

	Context        context.Context
	Environment    x.Services
	CustomContext  config.Context
	Auth           *auth0.Settings
	CachedToken    *token.CachedToken
	DecisionLogger *decisionlogger.Settings

	UI *clui.UI
}

func (ctx *CommonCtx) AccessToken() (string, error) {
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.Access, nil
}

func (ctx *CommonCtx) Token() (*api.Token, error) {
	return ctx.CachedToken.Get()
}

func (ctx *CommonCtx) AuthorizerAPIKey() (string, error) {
	tenantID := ctx.TenantID()
	kr, err := keyring.NewTenantKeyRing(tenantID)
	if err != nil {
		log.Printf("token: instantiating keyring, %s", err.Error())
		return "", nil
	}

	tkn, err := kr.GetToken()
	if err != nil {
		return "", nil
	}

	return tkn.AuthorizerAPIKey, nil
}

func (ctx *CommonCtx) DecisionLogsKey() (string, error) {
	tenantID := ctx.TenantID()
	kr, err := keyring.NewTenantKeyRing(tenantID)
	if err != nil {
		log.Printf("token: instantiating keyring, %s", err.Error())
		return "", nil
	}

	tkn, err := kr.GetToken()
	if err != nil {
		return "", nil
	}

	return tkn.DecisionLogsKey, nil
}

func (ctx *CommonCtx) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
