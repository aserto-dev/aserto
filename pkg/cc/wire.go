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
	decisionlogger "github.com/aserto-dev/aserto/pkg/decision_logger"
	"github.com/google/wire"
)

var (
	commonSet = wire.NewSet(
		iostream.NewUI,
		GetCacheKey,
		token.Load,
		NewAuthSettings,
		decisionlogger.NewSettings,
		clients.NewClientFactory,

		wire.Bind(new(clients.Factory), new(*clients.AsertoFactory)),
		wire.FieldsOf(new(*config.Config), "Services", "Context", "Auth", "DecisionLogger"),
		wire.Struct(new(CommonCtx), "*"),
	)

	ccSet = wire.NewSet(
		config.NewConfig,
		commonSet,

		iostream.DefaultIO,
		context.Background,

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
	overrides ...config.Overrider,
) (*CommonCtx, error) {
	wire.Build(ccSet)
	return &CommonCtx{}, nil
}

func BuildTestCtx(
	ioStreams iostream.IO,
	configReader io.Reader,
	overrides ...config.Overrider,
) (*CommonCtx, error) {
	wire.Build(ccTestSet)
	return &CommonCtx{}, nil
}

func GetCacheKey(auth *config.Auth) token.CacheKey {
	return token.CacheKey(auth.Issuer)
}

func NewAuthSettings(auth *config.Auth) *auth0.Settings {
	return auth.GetSettings()
}
