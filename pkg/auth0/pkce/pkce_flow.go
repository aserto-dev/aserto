// Package pkce implements the OAuth PKCE Authorization Flow for client applications by
// starting a server at localhost to receive the web redirect after the user has authorized the application.
package pkce

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
)

// Flow holds the state for the steps of OAuth Web Application flow.
type Flow struct {
	server   *localServer
	clientID string
	state    string
}

// InitFlow creates a new Flow instance by detecting a locally available port number.
func InitFlow(addr string) (*Flow, error) {
	server, err := bindLocalServer(addr)
	if err != nil {
		return nil, err
	}

	state, _ := randomString(20)

	return &Flow{
		server: server,
		state:  state,
	}, nil
}

// BrowserParams are GET query parameters for initiating the PKCE flow.
type BrowserParams struct {
	ClientID      string
	Audience      string
	RedirectURI   string
	Scopes        []string
	CodeVerifier  string
	CodeChallenge string
}

// BrowserURL appends GET query parameters to baseURL and returns the url that the user should
// navigate to in their web browser.
// nolint:gocritic // external code too much risk changing
func (flow *Flow) BrowserURL(baseURL string, params BrowserParams) (string, error) {
	ru, err := url.Parse(params.RedirectURI)
	if err != nil {
		return "", err
	}

	ru.Host = fmt.Sprintf("%s:%s", ru.Hostname(), ru.Port())
	flow.server.CallbackPath = ru.Path
	flow.clientID = params.ClientID

	q := url.Values{}
	q.Set("scope", strings.Join(params.Scopes, " "))
	q.Set("response_type", "code")
	q.Set("client_id", params.ClientID)
	q.Set("code_challenge", params.CodeChallenge)
	q.Set("code_challenge_method", "S256")
	q.Set("redirect_uri", ru.String())
	q.Set("state", flow.state)
	q.Set("audience", params.Audience)
	q.Set("prompt", "login")

	return fmt.Sprintf("%s?%s", baseURL, q.Encode()), nil
}

// StartServer starts the localhost server and blocks until it has received the web redirect. The
// writeSuccess function can be used to render a HTML page to the user upon completion.
func (flow *Flow) StartServer(writeSuccess func(io.Writer)) error {
	flow.server.WriteSuccessHTML = writeSuccess
	return flow.server.Serve()
}

var errStateMismatch = errors.New("state mismatch")

// AccessToken blocks until the browser flow has completed and returns the access token.
func (flow *Flow) AccessToken(ctx context.Context, tokenURL, redirectURL, codeVerifier string) (*api.Token, error) {
	code, err := flow.server.WaitForCode()
	if err != nil {
		return nil, err
	}
	if code.State != flow.state {
		return nil, errStateMismatch
	}

	resp, err := api.PostForm(ctx,
		tokenURL,
		url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {flow.clientID},
			"code_verifier": {codeVerifier},
			"code":          {code.Code},
			"redirect_uri":  {redirectURL},
		})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func randomString(length int) (string, error) {
	b := make([]byte, length/2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type CodeChallenge struct {
	Verifier  string
	Challenge string
}

var errInvalidArg = errors.New("invalid argument")

func CreateCodeChallenge(n int) (CodeChallenge, error) {
	const safe = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-._~"
	if n < 43 || n > 128 {
		return CodeChallenge{"", ""}, errInvalidArg
	}
	buff := make([]byte, n)
	if _, err := rand.Read(buff); err != nil {
		return CodeChallenge{"", ""}, err
	}
	nsafe := byte(len(safe))
	for i, b := range buff {
		b %= nsafe
		buff[i] = safe[b]
	}
	cv := base64.RawURLEncoding.EncodeToString(buff)
	s256 := sha256.Sum256([]byte(cv))
	cc := base64.RawURLEncoding.EncodeToString(s256[:])
	return CodeChallenge{cv, cc}, nil
}
