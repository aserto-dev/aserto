package cc

import (
	"context"
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
	decisionlogger "github.com/aserto-dev/aserto/pkg/decision_logger"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/aserto/pkg/x"
	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
)

const (
	TenantSuffix string = ".aserto.com"
)

type CommonCtx struct {
	clients.Factory

	Config         *config.Config
	Context        context.Context
	Environment    *x.Services
	Auth           *auth0.Settings
	CachedToken    *token.CachedToken
	TopazContext   *topazCC.CommonCtx
	DecisionLogger *decisionlogger.Settings
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
