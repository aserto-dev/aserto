package cc

import (
	"context"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	dl "github.com/aserto-dev/aserto/pkg/decision_logger"
)

// NewCommonCtx, CommonContext constructor (extracted from wire).
func NewCommonCtx(configPath config.Path, overrides ...config.Overrider) (*CommonCtx, error) {
	contextContext := context.Background()
	configConfig, err := config.NewConfig(configPath, overrides...)
	if err != nil {
		return nil, err
	}

	services := &configConfig.Services
	auth := configConfig.Auth

	cacheKey := GetCacheKey(auth)
	cachedToken := token.Load(cacheKey)

	tenantID := NewTenantID(configConfig, cachedToken)
	asertoFactory, err := clients.NewClientFactory(contextContext, services, tenantID, cachedToken)
	if err != nil {
		return nil, err
	}

	settings := NewAuthSettings(auth)

	dlConfig := &configConfig.DecisionLogger
	dlSettings := dl.NewSettings(dlConfig)

	commonCtx := &CommonCtx{
		Factory:        asertoFactory,
		Context:        contextContext,
		Config:         configConfig,
		Environment:    services,
		Auth:           settings,
		CachedToken:    cachedToken,
		DecisionLogger: dlSettings,
	}

	return commonCtx, nil
}

func NewTenantID(cfg *config.Config, cachedToken *token.CachedToken) clients.TenantID {
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
