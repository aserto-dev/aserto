package auth0

import "github.com/aserto-dev/aserto/pkg/x"

type Settings struct {
	Environment            string
	Domain                 string
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

func GetSettings(env string) *Settings {

	switch env {
	case x.EnvProduction:
		return &Settings{
			Environment:            x.EnvProduction,
			Domain:                 "aserto.us.auth0.com",
			Audience:               "https://console.aserto.com",
			ClientID:               "98ofxNoUdgVu7vuYAddWW2WpglFM4til",
			RedirectURL:            "http://localhost:3987",
			LogoutURL:              "http://localhost:3987",
			AuthorizationURL:       "https://aserto.us.auth0.com/authorize",
			DeviceAuthorizationURL: "https://aserto.us.auth0.com/oauth/device/code",
			TokenURL:               "https://aserto.us.auth0.com/oauth/token",
			UserInfoURL:            "https://aserto.us.auth0.com/userinfo",
			OpenIDConfiguration:    "https://aserto.us.auth0.com/.well-known/openid-configuration",
			JWKS:                   "https://aserto.us.auth0.com/.well-known/jwks.json",
		}
	case x.EnvEngineering:
		return &Settings{
			Environment:            x.EnvEngineering,
			Domain:                 "aserto-eng.us.auth0.com",
			Audience:               "https://console.aserto.com",
			ClientID:               "IZFpro8wrS35QjSjAbQj7ylBqjutNdXe",
			RedirectURL:            "http://localhost:3987",
			LogoutURL:              "http://localhost:3987",
			AuthorizationURL:       "https://aserto-eng.us.auth0.com/authorize",
			DeviceAuthorizationURL: "https://aserto-eng.us.auth0.com/oauth/device/code",
			TokenURL:               "https://aserto-eng.us.auth0.com/oauth/token",
			UserInfoURL:            "https://aserto-eng.us.auth0.com/userinfo",
			OpenIDConfiguration:    "https://aserto-eng.us.auth0.com/.well-known/openid-configuration",
			JWKS:                   "https://aserto-eng.us.auth0.com/.well-known/jwks.json",
		}
	default:
		return nil
	}
}
