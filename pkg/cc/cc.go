package cc

import (
	"context"
	"log"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/aserto-dev/clui"

	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
)

type CommonCtx struct {
	clients.Factory

	Context       context.Context
	Environment   x.Services
	CustomContext config.Context
	Auth          *auth0.Settings
	CachedToken   *token.CachedToken
	TopazContext  *topazCC.CommonCtx
	UI            *clui.UI
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
	cachedTkn, err := ctx.CachedToken.Get()
	if err != nil {
		log.Printf("token: failed to retrieve cached token, %s", err.Error())
		return "", err
	}

	kr, err := keyring.NewTenantKeyRing(tenantID + "-" + cachedTkn.Subject)
	if err != nil {
		log.Printf("token: instantiating keyring, %s", err.Error())
		return "", err
	}

	tkn, err := kr.GetToken()
	if err != nil {
		return "", err
	}

	return tkn.AuthorizerAPIKey, nil
}

func (ctx *CommonCtx) DecisionLogsKey() (string, error) {
	tenantID := ctx.TenantID()
	cachedTkn, err := ctx.CachedToken.Get()
	if err != nil {
		log.Printf("token: failed to retrieve cached token, %s", err.Error())
		return "", nil
	}

	kr, err := keyring.NewTenantKeyRing(tenantID + "-" + cachedTkn.Subject)
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
