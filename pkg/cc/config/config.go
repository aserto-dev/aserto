package config

import (
	"strings"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	EnvironmentErr = errors.New("unknown environment")
)

// Overrider is a func that mutates configuration
type Overrider func(*Config)

type Auth struct {
	Issuer   string `json:"issuer"`
	ClientID string `json:"client_id"`
	Audience string `json:"audience"`
}

func (auth *Auth) GetSettings() *auth0.Settings {
	return auth0.GetSettings(auth.Issuer, auth.ClientID, auth.Audience)
}

type Config struct {
	TenantID string     `json:"tenant_id"`
	Services x.Services `json:"services"`
	Auth     *Auth      `json:"auth"`
}

type Path string

func NewConfig(path Path, overrides Overrider) (*Config, error) {
	cfg := new(Config)

	v := viper.New()
	v.SetEnvPrefix("ASERTO")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.SetDefault("tenant_id", "")
	v.SetDefault("services", x.DefaultEnvironment())
	v.SetDefault("auth.issuer", auth0.IssuerProduction)
	v.SetDefault("auth.client_id", auth0.ClientIDProduction)
	v.SetDefault("auth.audience", auth0.Audience)

	configFile := string(path)
	if path != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, errors.Wrapf(err, "failed to read config file [%s]", configFile)
		}
	}

	v.AutomaticEnv()

	if err := v.UnmarshalExact(cfg, jsonDecoderConfig); err != nil {
		return nil, errors.Wrap(err, "failed to parse config file")
	}

	if overrides != nil {
		overrides(cfg)
	}

	return cfg, nil
}

func jsonDecoderConfig(dc *mapstructure.DecoderConfig) {
	dc.TagName = "json"
}
