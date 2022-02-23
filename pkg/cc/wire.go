//go:build wireinject
// +build wireinject

package cc

import (
	"context"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	"github.com/aserto-dev/clui"
	"github.com/google/wire"
)

var (
	ccSet = wire.NewSet(
		context.Background,
		config.NewConfig,

		GetCacheKey,
		token.NewCachedToken,
		NewTenantID,
		NewAuthSettings,
		clients.NewClientFactory,

		clui.NewUI,

		wire.Bind(new(clients.Factory), new(*clients.AsertoFactory)),
		wire.FieldsOf(new(*config.Config), "Services", "Auth"),
		wire.Struct(new(CommonCtx), "*"),
	)
)

func BuildCommonCtx(
	configPath config.Path,
	overrides config.Overrider,
	svcOptions *clients.ServiceOptions,
) (*CommonCtx, error) {
	wire.Build(ccSet)
	return &CommonCtx{}, nil
}

func NewTenantID(cfg *config.Config, cachedToken token.CachedToken) clients.TenantID {
	id := cfg.TenantID
	if id == "" {
		id = cachedToken.TenantID()
	}

	return clients.TenantID(id)
}

func GetCacheKey(auth *config.Auth) token.CacheKey {
	return token.CacheKey(auth.Issuer)
}

func NewAuthSettings(auth *config.Auth) *auth0.Settings {
	return auth.GetSettings()
}