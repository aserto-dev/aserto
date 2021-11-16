package auth0

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const (
	envAuth0Domain       = "AUTH0_DOMAIN"
	envAuth0ClientID     = "AUTH0_CLIENT_ID"
	envAuth0ClientSecret = "AUTH0_CLIENT_SECRET" // nolint:gosec // not a hardcoded credential
)

// Config -.
type Config struct {
	Domain       string `mapstructure:"domain"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// FromEnv - create config instance from environment variables.
func FromEnv() (*Config, error) {
	a := Config{
		Domain:       os.Getenv(envAuth0Domain),
		ClientID:     os.Getenv(envAuth0ClientID),
		ClientSecret: os.Getenv(envAuth0ClientSecret),
	}
	return &a, nil
}

func FromProfile(profile string) (*Config, error) {
	r, err := os.Open(profile)
	if err != nil {
		return nil, errors.Wrapf(err, "opening profile [%s]", profile)
	}

	envMap, err := godotenv.Parse(r)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing profile [%s]", profile)
	}

	a := Config{}
	var ok bool
	if a.Domain, ok = envMap[envAuth0Domain]; !ok {
		return nil, errors.Errorf("[%s] element is not specified in profile [%s]", envAuth0Domain, profile)
	}
	if a.ClientID, ok = envMap[envAuth0ClientID]; !ok {
		return nil, errors.Errorf("[%s] element is not specified in profile [%s]", envAuth0ClientID, profile)
	}
	if a.ClientSecret, ok = envMap[envAuth0ClientSecret]; !ok {
		return nil, errors.Errorf("[%s] element is not specified in profile [%s]", envAuth0ClientSecret, profile)
	}

	return &a, nil
}

func (a *Config) Validate() error {
	if a == nil {
		return errors.Errorf("auth0 config not initialized")
	}

	switch {
	case a.Domain == "":
		return errors.Errorf("auth0 domain configuration setting is missing")
	case a.ClientID == "":
		return errors.Errorf("auth0 client id configuration setting is missing")
	case a.ClientSecret == "":
		return errors.Errorf("auth0 client secret configuration setting is missing")
	}
	return nil
}
