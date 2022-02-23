package user

import (
	"context"
	"fmt"
	"net/url"
	"time"

	aserto "github.com/aserto-dev/aserto-go/client"
	"github.com/aserto-dev/aserto-go/client/tenant"
	auth0api "github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/aserto-dev/aserto/pkg/auth0/pkce"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/keyring"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	account "github.com/aserto-dev/go-grpc/aserto/tenant/account/v1"
	connection "github.com/aserto-dev/go-grpc/aserto/tenant/connection/v1"

	"github.com/cli/browser"
	"github.com/pkg/errors"
)

type LoginCmd struct{}

func (d *LoginCmd) Run(c *cc.CommonCtx) error {
	settings := c.Auth
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

	codeChallenge, err := pkce.CreateCodeChallenge(50)
	if err != nil {
		return err
	}

	params := pkce.BrowserParams{
		ClientID:      settings.ClientID,
		RedirectURI:   settings.RedirectURL,
		Audience:      settings.Audience,
		Scopes:        scopes,
		CodeVerifier:  codeChallenge.Verifier,
		CodeChallenge: codeChallenge.Challenge,
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

	tok, err := flow.AccessToken(c.Context, settings.TokenURL, settings.RedirectURL, codeChallenge.Verifier)
	if err != nil {
		return err
	}

	tok.ExpiresAt = time.Now().UTC().Add(time.Second * time.Duration(tok.ExpiresIn))

	conn, err := tenant.New(
		c.Context,
		aserto.WithAddr(c.Environment.TenantService.Address),
		aserto.WithTokenAuth(tok.Access),
	)
	if err != nil {
		return err
	}

	if err = getTenantID(c.Context, conn, tok); err != nil {
		return errors.Wrapf(err, "get tenant id")
	}

	if err = GetConnectionKeys(c.Context, conn, tok); err != nil {
		return errors.Wrapf(err, "get connection keys")
	}

	kr, err := keyring.NewKeyRing(c.Auth.Issuer)
	if err != nil {
		return err
	}
	if err := kr.SetToken(tok); err != nil {
		return err
	}

	fmt.Fprint(c.UI.Err(), "login successful\n")

	return nil
}

func getTenantID(ctx context.Context, client *tenant.Client, tok *auth0api.Token) error {
	resp, err := client.Account.GetAccount(ctx, &account.GetAccountRequest{})
	if err != nil {
		return errors.Wrapf(err, "get account")
	}

	tok.TenantID = resp.Result.DefaultTenant

	return err
}

func GetConnectionKeys(ctx context.Context, client *tenant.Client, tok *auth0api.Token) error {
	client.SetTenantID(tok.TenantID)

	resp, err := client.Connections.ListConnections(
		ctx,
		&connection.ListConnectionsRequest{
			Kind: api.ProviderKind_PROVIDER_KIND_UNKNOWN,
		})
	if err != nil {
		return errors.Wrapf(err, "list connections account")
	}

	//nolint:exhaustive // we only care about these two provider kinds.
	for _, cn := range resp.Results {
		switch cn.Kind {
		case api.ProviderKind_PROVIDER_KIND_AUTHORIZER:
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				tok.AuthorizerAPIKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get connection [%s]", cn.Id)
			}
		case api.ProviderKind_PROVIDER_KIND_DECISION_LOGS:
			if respX, err := GetConnection(ctx, client, cn.Id); err == nil {
				tok.DecisionLogsKey = respX.Result.Config.Fields["api_key"].GetStringValue()
			} else {
				return errors.Wrapf(err, "get connection [%s]", cn.Id)
			}
		}
	}

	return nil
}

func GetConnection(
	ctx context.Context,
	client *tenant.Client,
	connectionID string,
) (*connection.GetConnectionResponse, error) {
	return client.Connections.GetConnection(
		ctx,
		&connection.GetConnectionRequest{Id: connectionID},
	)
}
