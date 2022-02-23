// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func BuildCommonCtx(configPath config.Path, overrides config.Overrider, svcOptions *clients.ServiceOptions) (*CommonCtx, error) {
	contextContext := context.Background()
	configConfig, err := config.NewConfig(configPath, overrides)
	if err != nil {
		return nil, err
	}
	services := &configConfig.Services
	auth := configConfig.Auth
	cacheKey := GetCacheKey(auth)
	cachedToken := token.NewCachedToken(cacheKey)
	tenantID := NewTenantID(configConfig, cachedToken)
	asertoFactory, err := clients.NewClientFactory(contextContext, svcOptions, services, tenantID, cachedToken)
	if err != nil {
		return nil, err
	}
	settings := NewAuthSettings(auth)
	ui := clui.NewUI()
	commonCtx := &CommonCtx{
		Factory:     asertoFactory,
		Context:     contextContext,
		Environment: services,
		Auth:        settings,
		CachedToken: cachedToken,
		UI:          ui,
	}
	return commonCtx, nil
}

// wire.go:

var (
	ccSet = wire.NewSet(context.Background, config.NewConfig, GetCacheKey, token.NewCachedToken, NewTenantID,
		NewAuthSettings, clients.NewClientFactory, clui.NewUI, wire.Bind(new(clients.Factory), new(*clients.AsertoFactory)), wire.FieldsOf(new(*config.Config), "Services", "Auth"), wire.Struct(new(CommonCtx), "*"),
	)
)

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
