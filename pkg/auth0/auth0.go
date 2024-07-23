package auth0

import "fmt"

const (
	Issuer    string = "aserto.us.auth0.com"
	ClientID  string = "wxB6q804bWWiPqWRtauOeBGtZobfBWD9"
	Audience  string = "https://cli.aserto.com"
	GrantType string = "urn:ietf:params:oauth:grant-type:device_code"
)

type Settings struct {
	Issuer                 string
	Audience               string
	ClientID               string
	GrantType              string
	RedirectURL            string
	LogoutURL              string
	AuthorizationURL       string
	DeviceAuthorizationURL string
	TokenURL               string
	UserInfoURL            string
	OpenIDConfiguration    string
	JWKS                   string
}

func GetSettings(issuer, clientID, audience string) *Settings {
	return &Settings{
		Issuer:                 issuer,
		Audience:               audience,
		ClientID:               clientID,
		RedirectURL:            "http://localhost:3987",
		LogoutURL:              "http://localhost:3987",
		GrantType:              GrantType,
		AuthorizationURL:       fmt.Sprintf("https://%s/authorize", issuer),
		DeviceAuthorizationURL: fmt.Sprintf("https://%s/oauth/device/code", issuer),
		TokenURL:               fmt.Sprintf("https://%s/oauth/token", issuer),
		UserInfoURL:            fmt.Sprintf("https://%s/userinfo", issuer),
		OpenIDConfiguration:    fmt.Sprintf("https://%s/.well-known/openid-configuration", issuer),
		JWKS:                   fmt.Sprintf("https://%s/.well-known/jwks.json", issuer),
	}
}
