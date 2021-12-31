package cc

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/aserto/pkg/paths"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

type CommonCtx struct {
	Context     context.Context
	OutWriter   io.Writer
	ErrWriter   io.Writer
	Insecure    bool
	environment string
	services    *x.Services
	_token      *api.Token
	overrides   map[string]string
}

func New() *CommonCtx {
	log.SetOutput(io.Discard)
	log.SetPrefix("")
	log.SetFlags(log.LstdFlags)
	return &CommonCtx{
		Context:   context.Background(),
		OutWriter: os.Stdout,
		ErrWriter: os.Stderr,
		overrides: make(map[string]string),
	}
}

func (ctx *CommonCtx) SetEnv(env string) error {
	log.Printf("set-context-env %s", env)
	if env == "" {
		return errors.Errorf("env is not set")
	}

	var err error
	ctx.services, err = x.Environment(env)
	if err != nil {
		return err
	}

	ctx.environment = env

	return nil
}

func (ctx *CommonCtx) Environment() string {
	log.Printf("get-context-env %s", ctx.environment)
	return ctx.environment
}

func (ctx *CommonCtx) Override(key, value string) {
	log.Println("override-context-env", key, value)
	ctx.overrides[key] = value
}

func (ctx *CommonCtx) VerifyLoggedIn() error {
	if !ctx.isLoggedIn() {
		return errors.Errorf("user is not logged in, please login using '%s login'", x.AppName)
	}
	if ctx.isExpired() {
		return errors.Errorf("the access token has expired, please login using '%s login'", x.AppName)
	}
	return nil
}

func (ctx *CommonCtx) isLoggedIn() bool {
	if ctx.token() == nil || ctx.token().Access == "" {
		return false
	}
	return true
}

func (ctx *CommonCtx) isExpired() bool {
	return time.Now().UTC().After(ctx.token().ExpiresAt)
}

func (ctx *CommonCtx) AccessToken() string {
	return ctx.token().Access
}

func (ctx *CommonCtx) Token() *api.Token {
	return ctx.token()
}

func (ctx *CommonCtx) ExpiresAt() time.Time {
	return ctx.token().ExpiresAt
}

func (ctx *CommonCtx) TenantID() string {
	if tenantID, ok := ctx.overrides[x.TenantIDOverride]; ok {
		fmt.Fprintf(ctx.ErrWriter, "!!! tenant override [%s]\n", tenantID)
		return tenantID
	}
	return ctx.token().TenantID
}

func (ctx *CommonCtx) TenantSvcConnectionOptions(opts ...aserto.ConnectionOption) []aserto.ConnectionOption {
	return ctx.SvcConnectionOptions(ctx.TenantService(), opts...)
}

func (ctx *CommonCtx) AuthorizerSvcConnectionOptions(opts ...aserto.ConnectionOption) []aserto.ConnectionOption {
	return ctx.SvcConnectionOptions(ctx.AuthorizerService(), opts...)
}

func (ctx *CommonCtx) SvcConnectionOptions(addr string, opts ...aserto.ConnectionOption) []aserto.ConnectionOption {
	options := []aserto.ConnectionOption{
		aserto.WithAddr(addr),
		aserto.WithTokenAuth(ctx.AccessToken()),
		aserto.WithTenantID(ctx.TenantID()),
		aserto.WithInsecure(ctx.Insecure),
	}

	if strings.Contains(addr, "localhost") && !ctx.Insecure {
		p, err := paths.New()
		if err == nil {
			options = append(options, aserto.WithCACertPath(p.Certs.GRPC.CA))
		} else {
			fmt.Fprintln(ctx.ErrWriter, "Unable to locate onebox certificates.", err.Error())
		}
	}

	return append(options, opts...)
}

func (ctx *CommonCtx) AuthorizerService() string {
	if authorizer, ok := ctx.overrides[x.AuthorizerOverride]; ok {
		fmt.Fprintf(ctx.ErrWriter, "!!! authorizer override [%s]\n", authorizer)
		return authorizer
	}
	return ctx.services.AuthorizerService
}

func (ctx *CommonCtx) AuthorizerAPIKey() string {
	return ctx.token().AuthorizerAPIKey
}

func (ctx *CommonCtx) RegistryDownloadKey() string {
	return ctx.token().RegistryDownloadKey
}

func (ctx *CommonCtx) RegistryUploadKey() string {
	return ctx.token().RegistryUploadKey
}

func (ctx *CommonCtx) DecisionLogsKey() string {
	return ctx.token().DecisionLogsKey
}

func (ctx *CommonCtx) TenantService() string {
	return ctx.services.TenantService
}

func (ctx *CommonCtx) TasksService() string {
	return ctx.services.TasksService
}

func (ctx *CommonCtx) DecisionLogsService() string {
	return ctx.services.DecisionLogsService
}

func (ctx *CommonCtx) token() *api.Token {
	if ctx._token == nil {
		kr, err := keyring.NewKeyRing(ctx.environment)
		if err != nil {
			log.Printf("token: instantiating keyring, %s", err.Error())
			return nil
		}

		ctx._token, err = kr.GetToken()
		if err != nil {
			return nil
		}
	}
	return ctx._token
}

func (ctx *CommonCtx) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (ctx *CommonCtx) SetLogger(w io.Writer) {
	log.SetOutput(w)
}
