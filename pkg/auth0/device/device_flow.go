package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	u "net/url"
	"os"
	"strings"
	"time"

	"github.com/aserto-dev/aserto/pkg/auth0/api"
	"github.com/lestrrat-go/jwx/jwt"
)

type DeviceCodeFlow struct {
	DeviceAuthorizationURL string
	TokenURL               string
	ClientID               string
	Audience               string
	GrantType              string
	Scopes                 []string
	deviceCode             *DeviceCode
	accessToken            *TokenResponse
}

type DeviceCodeOption func(*DeviceCodeFlow)

func New(opts ...DeviceCodeOption) *DeviceCodeFlow {
	d := &DeviceCodeFlow{}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func WithDeviceAuthorizationURL(url string) DeviceCodeOption {
	return func(i *DeviceCodeFlow) {
		i.DeviceAuthorizationURL = url
	}
}

func WithTokenURL(url string) DeviceCodeOption {
	return func(i *DeviceCodeFlow) {
		i.TokenURL = url
	}
}

func WithClientID(id string) DeviceCodeOption {
	return func(i *DeviceCodeFlow) {
		i.ClientID = id
	}
}

func WithAudience(audience string) DeviceCodeOption {
	return func(i *DeviceCodeFlow) {
		i.Audience = audience
	}
}

func WithGrantType(grantType string) DeviceCodeOption {
	return func(i *DeviceCodeFlow) {
		i.GrantType = grantType
	}
}

func WithScope(scopes ...string) DeviceCodeOption {
	return func(i *DeviceCodeFlow) {
		i.Scopes = append(i.Scopes, scopes...)
	}
}

func (f *DeviceCodeFlow) Reader() io.Reader {
	q := u.Values{}
	q.Set("client_id", f.ClientID)
	q.Set("scope", strings.Join(f.Scopes, " "))
	q.Set("audience", f.Audience)
	return strings.NewReader(q.Encode())
}

type DeviceCode struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type TokenRequest struct {
	URL        string
	GrantType  string
	DeviceCode string
	ClientID   string
	ExpiresIn  int
	Interval   int
}

func (f *DeviceCodeFlow) TokenReader() io.Reader {
	q := u.Values{}
	q.Set("grant_type", f.GrantType)
	q.Set("device_code", f.deviceCode.DeviceCode)
	q.Set("client_id", f.ClientID)
	return strings.NewReader(q.Encode())
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	StatusCode       int    `json:"status_code"`
}

func (f *DeviceCodeFlow) GetDeviceCode(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "POST", f.DeviceAuthorizationURL, f.Reader())
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var resp DeviceCode

	if res.StatusCode == http.StatusOK {
		if err := json.Unmarshal(body, &resp); err != nil {
			return err
		}
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Status %s %d\n", res.Status, res.StatusCode)
		fmt.Println(res)
	}

	f.deviceCode = &resp

	return nil
}

func (f *DeviceCodeFlow) RequestAccessToken(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", f.TokenURL, f.TokenReader())
	if err != nil {
		return false, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	var resp TokenResponse
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return false, err
	}

	f.accessToken = &resp

	return res.StatusCode == http.StatusOK, nil
}

func (f *DeviceCodeFlow) AccessToken() *api.Token {
	if f.accessToken == nil {
		return nil
	}

	options := []jwt.ParseOption{
		jwt.WithValidate(true),
		jwt.WithAcceptableSkew(time.Duration(2) * time.Second),
	}

	jwtToken, err := jwt.ParseString(
		f.accessToken.AccessToken,
		options...,
	)
	if err != nil {
		return nil
	}

	subjectRunes := strings.Split(jwtToken.Subject(), "|")

	var sub string
	if len(subjectRunes) == 2 {
		sub = subjectRunes[1]
	} else {
		sub = jwtToken.Subject()
	}

	return &api.Token{
		Type:      f.accessToken.TokenType,
		Scope:     strings.Join(f.Scopes, " "),
		Identity:  f.accessToken.IDToken,
		Access:    f.accessToken.AccessToken,
		Subject:   sub,
		ExpiresIn: f.accessToken.ExpiresIn,
		ExpiresAt: time.Now().UTC().Add(time.Second * time.Duration(f.accessToken.ExpiresIn)),
	}
}

func (f *DeviceCodeFlow) GetUserCode() string {
	if f.deviceCode == nil {
		return ""
	}
	return f.deviceCode.UserCode
}

func (f *DeviceCodeFlow) GetVerificationURI() string {
	if f.deviceCode == nil {
		return ""
	}
	return f.deviceCode.VerificationURI
}

func (f *DeviceCodeFlow) GetVerificationURIComplete() string {
	if f.deviceCode == nil {
		return ""
	}
	return f.deviceCode.VerificationURIComplete
}

func (f *DeviceCodeFlow) ExpiresIn() time.Duration {
	if f.deviceCode == nil {
		return 0
	}
	return time.Duration(f.deviceCode.ExpiresIn) * time.Second
}

func (f *DeviceCodeFlow) Interval() time.Duration {
	if f.deviceCode == nil {
		return 0
	}
	return time.Duration(f.deviceCode.Interval) * time.Second
}
