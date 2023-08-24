package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aserto-dev/aserto/pkg/auth0"
	aErr "github.com/aserto-dev/aserto/pkg/cc/errors"
	decisionlogger "github.com/aserto-dev/aserto/pkg/decision_logger"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Overrider is a func that mutates configuration.
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
	Context        Context               `json:"context"`
	Services       x.Services            `json:"services"`
	Auth           *Auth                 `json:"auth"`
	DecisionLogger decisionlogger.Config `json:"decision_logger"`
}

type Context struct {
	Contexts      []Ctx  `json:"contexts"`
	ActiveContext string `json:"active"`
}

type Ctx struct {
	Name              string           `json:"name"`
	TenantID          string           `json:"tenant_id"`
	AuthorizerService x.ServiceOptions `json:"authorizer"`
}

type Path string

func NewConfig(path Path, overrides ...Overrider) (*Config, error) {
	configFile := string(path)

	return newConfig(
		func(v *viper.Viper) error {
			if configFile != "" {
				v.SetConfigFile(configFile)
				if err := v.ReadInConfig(); err != nil {
					return errors.Wrapf(err, "failed to read config file [%s]", configFile)
				}
			} else {
				filePath, err := GetConfigFile()
				if err != nil {
					return err
				}
				if _, err := os.Stat(filePath); err != nil {
					// user is not logged in yet
					return nil
				}
				v.SetConfigFile(filePath)
				if err := v.ReadInConfig(); err != nil {
					return errors.Wrapf(err, "failed to read config file [%s]", filePath)
				}
			}

			return nil
		},
		overrides...,
	)
}

func NewTestConfig(reader io.Reader, overrides ...Overrider) (*Config, error) {
	return newConfig(
		func(v *viper.Viper) error {
			v.SetConfigType("yaml")
			return v.ReadConfig(reader)
		},
		overrides...,
	)
}

type configReader func(*viper.Viper) error

func newConfig(reader configReader, overrides ...Overrider) (*Config, error) {
	cfg := new(Config)

	v := viper.New()
	v.SetEnvPrefix("ASERTO")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.SetDefault("services", x.DefaultEnvironment())
	v.SetDefault("auth.issuer", auth0.Issuer)
	v.SetDefault("auth.client_id", auth0.ClientID)
	v.SetDefault("auth.audience", auth0.Audience)

	v.AutomaticEnv()

	if err := reader(v); err != nil {
		return nil, err
	}

	if err := v.UnmarshalExact(cfg, jsonDecoderConfig); err != nil {
		return nil, errors.Wrap(err, "failed to parse config file")
	}

	for _, override := range overrides {
		override(cfg)
	}

	return cfg, nil
}

func jsonDecoderConfig(dc *mapstructure.DecoderConfig) {
	dc.TagName = "json"
}

func GetConfigFile() (string, error) {
	env := os.Getenv("ASERTO_ENV")
	if env != "" {
		return env, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine user home directory")
	}
	return filepath.Join(home, ".config", x.AppName, "config.json"), nil
}

func GetConfigFromFile(filePath string) (*Config, error) {
	if !FileExists(filePath) {
		return nil, aErr.NeedLoginErr
	}
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(dat, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
