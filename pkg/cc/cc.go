package cc

import (
	"context"
	"fmt"
	"io"
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
	ctx := CommonCtx{
		Context:   context.Background(),
		OutWriter: os.Stdout,
		ErrWriter: os.Stderr,
		overrides: make(map[string]string),
	}
	return &ctx
}

func (ctx *CommonCtx) SetEnv(env string) error {
	log.Printf("set-env %s", env)
	if env == "" {
		return errors.Errorf("env is not set")
	}

	ctx.services = grpcc.Environment(env)

	ctx.environment = env

	return nil
}

func (ctx *CommonCtx) Environment() string {
	return ctx.environment
}

func (ctx *CommonCtx) Override(key, value string) {
	log.Println("override", key, value)
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
