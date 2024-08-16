package cc

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/token"

	dl "github.com/aserto-dev/aserto/pkg/decision_logger"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/aserto/pkg/x"
	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
)

const (
	TenantSuffix string = ".aserto.com"
)

type CommonCtx struct {
	*topazCC.CommonCtx
	clients.Factory
	Config         *config.Config
	Environment    *x.Services
	Auth           *auth0.Settings
	CachedToken    *token.CachedToken
	DecisionLogger *dl.Settings
}

// NewCommonCtx, CommonContext constructor (extracted from wire).
func NewCommonCtx(tc *topazCC.CommonCtx, configPath config.Path, overrides ...config.Overrider) (*CommonCtx, error) {
	configConfig, err := config.NewConfig(configPath, overrides...)
	if err != nil {
		return nil, err
	}

	services := &configConfig.Services
	auth := configConfig.Auth

	cacheKey := GetCacheKey(auth)
	cachedToken := token.Load(cacheKey)

	configConfig.TenantID = newTenantID(configConfig, cachedToken)

	asertoFactory, err := clients.NewClientFactory(services, cachedToken)
	if err != nil {
		return nil, err
	}

	settings := newAuthSettings(auth)

	dlConfig := &configConfig.DecisionLogger
	dlSettings := dl.NewSettings(dlConfig)

	commonCtx := &CommonCtx{
		CommonCtx:      tc,
		Factory:        asertoFactory,
		Config:         configConfig,
		Environment:    services,
		Auth:           settings,
		CachedToken:    cachedToken,
		DecisionLogger: dlSettings,
	}

	return commonCtx, nil
}

func newTenantID(cfg *config.Config, cachedToken *token.CachedToken) string {
	id := cfg.TenantID
	if id == "" {
		id = cachedToken.TenantID()
	}
	return id
}

func GetCacheKey(auth *config.Auth) token.CacheKey {
	return token.CacheKey(auth.Issuer)
}

func newAuthSettings(auth *config.Auth) *auth0.Settings {
	return auth.GetSettings()
}

func (ctx *CommonCtx) TenantID() string {
	tkn, err := ctx.Token()
	if err != nil {
		return ""
	}
	return tkn.TenantID
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
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.AuthorizerAPIKey, nil
}

func (ctx *CommonCtx) DirectoryReadKey() (string, error) {
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.DirectoryReadKey, nil
}

func (ctx *CommonCtx) DirectoryWriteKey() (string, error) {
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.DirectoryWriteKey, nil
}

func (ctx *CommonCtx) DiscoveryKey() (string, error) {
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.DiscoveryKey, nil
}

func (ctx *CommonCtx) DecisionLogsKey() (string, error) {
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.DecisionLogsKey, nil
}

func (ctx *CommonCtx) RegistryReadKey() (string, error) {
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.RegistryDownloadKey, nil
}

func (ctx *CommonCtx) RegistryWriteKey() (string, error) {
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.RegistryUploadKey, nil
}

func (ctx *CommonCtx) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (ctx *CommonCtx) SaveContextConfig(configurationFile string) error {
	configDir := filepath.Dir(configurationFile)
	if !filex.DirExists(configDir) {
		err := os.MkdirAll(configDir, 0o700)
		if err != nil {
			return err
		}
	}
	kongConfigBytes, err := json.Marshal(ctx.Config)
	if err != nil {
		return err
	}
	err = os.WriteFile(configurationFile, kongConfigBytes, 0o600)
	if err != nil {
		return err
	}
	return nil
}

func IsAsertoAccount(name string) bool {
	return strings.HasSuffix(name, TenantSuffix)
}
