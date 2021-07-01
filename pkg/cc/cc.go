package cc

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/keyring"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

type CommonCtx struct {
	Context     context.Context
	OutWriter   io.Writer
	ErrWriter   io.Writer
	environment string
	services    *grpcc.Services
	_token      *api.Token
	overrides   map[string]string
}

func New() *CommonCtx {
	log.SetOutput(ioutil.Discard)
	log.SetPrefix("")
	log.SetFlags(log.LstdFlags)
	ctx := CommonCtx{
		Context:   context.Background(),
		OutWriter: os.Stdout,
		ErrWriter: os.Stderr,
		overrides: make(map[string]string),
	}
	return &ctx
}

func (ctx *CommonCtx) SetEnv(env string) error {
	log.Printf("set-context-env %s", env)
	if env == "" {
		return errors.Errorf("env is not set")
	}

	var err error
	ctx.services, err = grpcc.Environment(env)
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
	if !ctx.IsLoggedIn() {
		return errors.Errorf("user is not logged in, please login using '%s login'", x.AppName)
	}
	return nil
}

func (ctx *CommonCtx) IsLoggedIn() bool {
	if ctx.token() == nil || ctx.token().Access == "" {
		return false
	}
	return true
}

func (ctx *CommonCtx) AccessToken() string {
	return ctx.token().Access
}

func (ctx *CommonCtx) Token() *api.Token {
	return ctx.token()
}

func (ctx *CommonCtx) TenantID() string {
	if tenantID, ok := ctx.overrides[x.TenantIDOverride]; ok {
		fmt.Fprintf(ctx.ErrWriter, "!!! tenant override [%s]\n", tenantID)
		return tenantID
	}
	return ctx.token().TenantID
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

func (ctx *CommonCtx) RegistrySvc() string {
	return ctx.services.RegistryService
}

func (ctx *CommonCtx) RegistryDownloadKey() string {
	return ctx.token().RegistryDownloadKey
}

func (ctx *CommonCtx) RegistryUploadKey() string {
	return ctx.token().RegistryUploadKey
}

func (ctx *CommonCtx) TenantService() string {
	return ctx.services.TenantService
}

func (ctx *CommonCtx) TasksService() string {
	return ctx.services.TasksService
}

func (ctx *CommonCtx) token() *api.Token {
	if ctx._token == nil {
		kr, err := keyring.NewKeyRing()
		if err != nil {
			log.Printf("token: instantiating keyring, %s", err.Error())
			return nil
		}

		ctx._token, err = kr.GetToken(ctx.environment)
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
