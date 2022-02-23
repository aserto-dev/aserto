package auth0

import "fmt"

type Settings struct {
	Issuer                 string
	Audience               string
	ClientID               string
	RedirectURL            string
	LogoutURL              string
	AuthorizationURL       string
	DeviceAuthorizationURL string
	TokenURL               string
	UserInfoURL            string
	OpenIDConfiguration    string
	JWKS                   string
}

const (
	IssuerProduction   = "aserto.us.auth0.com"
	ClientIDProduction = "98ofxNoUdgVu7vuYAddWW2WpglFM4til"

	Audience = "https://console.aserto.com"
)

func GetSettings(issuer, clientID, audience string) *Settings {
	return &Settings{
		Issuer:                 issuer,
		Audience:               audience,
		ClientID:               clientID,
		RedirectURL:            "http://localhost:3987",
		LogoutURL:              "http://localhost:3987",
		AuthorizationURL:       fmt.Sprintf("https://%s/authorize", issuer),
		DeviceAuthorizationURL: fmt.Sprintf("https://%s/oauth/device/code", issuer),
		TokenURL:               fmt.Sprintf("https://%s/oauth/token", issuer),
		UserInfoURL:            fmt.Sprintf("https://%s/userinfo", issuer),
		OpenIDConfiguration:    fmt.Sprintf("https://%s/.well-known/openid-configuration", issuer),
		JWKS:                   fmt.Sprintf("https://%s/.well-known/jwks.json", issuer),
	}
}
