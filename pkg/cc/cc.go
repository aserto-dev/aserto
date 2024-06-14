package cc

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/cc/clients"
	"github.com/aserto-dev/aserto/pkg/cc/config"
	"github.com/aserto-dev/aserto/pkg/cc/token"
	decisionlogger "github.com/aserto-dev/aserto/pkg/decision_logger"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/aserto-dev/clui"
	topazCC "github.com/aserto-dev/topaz/pkg/cli/cc"
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
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.AuthorizerAPIKey, nil
}

func (ctx *CommonCtx) DecisionLogsKey() (string, error) {
	tkn, err := ctx.Token()
	if err != nil {
		return "", err
	}
	return tkn.DecisionLogsKey, nil
}

func (ctx *CommonCtx) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (ctx *CommonCtx) SaveContextConfig(configurationFile string) error {
	configDir := filepath.Dir(configurationFile)
	if !filex.DirExists(configDir) {
		err := os.MkdirAll(configDir, 0700)
		if err != nil {
			return err
		}
	}
	kongConfigBytes, err := json.Marshal(ctx.Config)
	if err != nil {
		return err
	}
	err = os.WriteFile(configurationFile, kongConfigBytes, 0666) // nolint
	if err != nil {
		return err
	}
	return nil
}

func IsAsertoAccount(name string) bool {
	isAsertoAccount, _ := regexp.MatchString(`\w+[.][0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, name)
	return isAsertoAccount
}
