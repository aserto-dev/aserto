//go:build wireinject
// +build wireinject

package cc

import (
	"context"
	"io"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/iostream"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	"github.com/google/wire"
)

var (
	commonSet = wire.NewSet(
		iostream.NewUI,
		GetCacheKey,
		token.Load,
		NewTenantID,
		NewAuthSettings,
		clients.NewClientFactory,

		wire.Bind(new(clients.Factory), new(*clients.AsertoFactory)),
		wire.FieldsOf(new(*config.Config), "Services", "Auth"),
		wire.Struct(new(CommonCtx), "*"),
	)

	ccSet = wire.NewSet(
		commonSet,

		iostream.DefaultIO,
		context.Background,
		config.NewConfig,

		wire.Bind(new(iostream.IO), new(*iostream.StdIO)),
	)

	ccTestSet = wire.NewSet(
		commonSet,

		context.TODO,
		config.NewTestConfig,
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

func BuildTestCtx(
	ioStreams iostream.IO,
	configReader io.Reader,
	overrides config.Overrider,
	svcOptions *clients.ServiceOptions,
) (*CommonCtx, error) {
	wire.Build(ccTestSet)
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
