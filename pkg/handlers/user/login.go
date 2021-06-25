package user

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/aserto-dev/aserto/pkg/auth0"
	auth0api "github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/auth0/pkce"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/grpcc"
	"github.com/aserto-dev/aserto/pkg/grpcc/tenant"
	"github.com/aserto-dev/aserto/pkg/keyring"

	"github.com/aserto-dev/proto/aserto/api"
	"github.com/aserto-dev/proto/aserto/tenant/account"
	"github.com/aserto-dev/proto/aserto/tenant/connection"

	"github.com/cli/browser"
	"github.com/pkg/errors"
)

type LoginCmd struct {
}

func (d *LoginCmd) Run(c *cc.CommonCtx) error {
	env := c.Environment()

	settings := auth0.GetSettings(env)
	scopes := []string{"openid", "email", "profile"}

	ru, err := url.Parse(settings.RedirectURL)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%s", ru.Hostname(), ru.Port())

	flow, err := pkce.InitFlow(addr)
	if err != nil {
		return err
	}

	codeVerifier, codeChallenge, err := pkce.CreateCodeChallenge(50)
	if err != nil {
		return err
	}

	params := pkce.BrowserParams{
		ClientID:      settings.ClientID,
		RedirectURI:   settings.RedirectURL,
		Audience:      settings.Audience,
		Scopes:        scopes,
		CodeVerifier:  codeVerifier,
		CodeChallenge: codeChallenge,
	}

	browserURL, err := flow.BrowserURL(settings.AuthorizationURL, params)
	if err != nil {
		return err
	}

	// A localhost server on a random available port will receive the web redirect.
	go func() {
		_ = flow.StartServer(nil)
	}()

	// Note: the user's web browser must run on the same device as the running app.
	err = browser.OpenURL(browserURL)
	if err != nil {
		return err
	}

	tok, err := flow.AccessToken(settings.TokenURL, settings.RedirectURL, codeVerifier)
	if err != nil {
		return err
	}

	tok.ExpiresAt = time.Now().UTC().Add(time.Second * time.Duration(tok.ExpiresIn))

	svcs, err := grpcc.Environment(env)
	if err != nil {
		return err
	}

	conn, err := tenant.Connection(
		c.Context,
		svcs.TenantService,
		grpcc.NewTokenAuth(tok.Access),
	)
	if err != nil {
		return err
	}

	if err := getTenantID(c.Context, conn, tok); err != nil {
		return errors.Wrapf(err, "get tenant id")
	}

	if err := getConnectionKeys(c.Context, conn, tok); err != nil {
		return errors.Wrapf(err, "get connection keys")
	}

	kr, err := keyring.NewKeyRing()
	if err != nil {
		return err
	}
	if err := kr.SetToken(env, tok); err != nil {
		return err
	}

	fmt.Fprint(c.ErrWriter, "login successful\n")

	return nil
}

func getTenantID(ctx context.Context, conn *tenant.Client, tok *auth0api.Token) error {
	accntClient := conn.AccountClient()

	resp, err := accntClient.GetAccount(ctx, &account.GetAccountRequest{})
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	tok.TenantID = resp.Result.DefaultTenant

	return err
}

func getConnectionKeys(ctx context.Context, conn *tenant.Client, tok *auth0api.Token) error {
	ctx = grpcc.SetTenantContext(ctx, tok.TenantID)

	connClient := conn.ConnectionManagerClient()
	resp, err := connClient.ListConnections(
		ctx,
		&connection.ListConnectionsRequest{
			Kind: api.ProviderKind_UNKNOWN_PROVIDER_KIND,
		})
	if err != nil {
		return errors.Wrapf(err, "list connections account")
	}

	//nolint:exhaustive // we only care about these two provider kinds.
	for _, cn := range resp.Results {
		switch cn.Kind {
		case api.ProviderKind_AUTHORIZER:
			respX, err := connClient.GetConnection(ctx, &connection.GetConnectionRequest{
				Id: cn.Id,
			})
			if err != nil {
				return errors.Wrapf(err, "get connection [%s]", cn.Id)
			}

			tok.AuthorizerAPIKey = respX.Result.Config.Fields["api_key"].GetStringValue()

		case api.ProviderKind_POLICY_REGISTRY:
			respX, err := connClient.GetConnection(ctx, &connection.GetConnectionRequest{
				Id: cn.Id,
			})
			if err != nil {
				return errors.Wrapf(err, "get connection [%s]", cn.Id)
			}

			tok.RegistryDownloadKey = respX.Result.Config.Fields["download_key"].GetStringValue()
			tok.RegistryUploadKey = respX.Result.Config.Fields["api_key"].GetStringValue()
		}
	}

	return err
}
