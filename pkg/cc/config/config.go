package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/filex"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const ConfigPath = "config.yaml"

// Overrider is a func that mutates configuration.
type Overrider func(*Config)

type Auth struct {
	Issuer   string `json:"issuer" yaml:"issuer"`
	ClientID string `json:"client_id" yaml:"client_id"`
	Audience string `json:"audience" yaml:"audience"`
}

func (auth *Auth) GetSettings() *auth0.Settings {
	return auth0.GetSettings(auth.Issuer, auth.ClientID, auth.Audience)
}

type Config struct {
	Context  Context    `json:"context" yaml:"context"`
	Services x.Services `json:"services" yaml:"services"`
	Auth     *Auth      `json:"auth" yaml:"auth"`
}

type Context struct {
	Contexts      []Ctx  `json:"contexts" yaml:"contexts"`
	ActiveContext string `json:"active" yaml:"active"`
}

type Ctx struct {
	Name              string            `json:"name" yaml:"name"`
	TenantID          string            `json:"tenant_id,omitempty" yaml:"tenant_id,omitempty"`
	AuthorizerService *x.ServiceOptions `json:"authorizer,omitempty" yaml:"authorizer,omitempty"`
	DirectoryReader   *x.ServiceOptions `json:"directory_reader,omitempty" yaml:"directory_reader,omitempty"`
	DirectoryWriter   *x.ServiceOptions `json:"directory_writer,omitempty" yaml:"directory_writer,omitempty"`
	DirectoryModel    *x.ServiceOptions `json:"directory_model,omitempty" yaml:"directory_model,omitempty"`
}

type Path string

func NewConfig(path Path, overrides ...Overrider) (*Config, error) {
	configFile := string(path)

	return newConfig(
		func(v *viper.Viper) error {
			if configFile != "" && filex.FileExists(configFile) {
				v.SetConfigFile(configFile)
				if err := v.ReadInConfig(); err != nil {
					return errors.Wrapf(err, "failed to read config file [%s]", configFile)
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
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetDefault("auth.issuer", auth0.Issuer)
	v.SetDefault("auth.client_id", auth0.ClientID)
	v.SetDefault("auth.audience", auth0.Audience)
	setDefaultServices(v, x.DefaultEnvironment())

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

func GetSymlinkConfigPath() (string, error) {
	envOverride := os.Getenv("ASERTO_DIR")
	if envOverride != "" {
		if filex.DirExists(envOverride) {
			return filepath.Join(envOverride, ConfigPath), nil
		}
		if filex.FileExists(envOverride) {
			return filepath.Join(filepath.Dir(envOverride), ConfigPath), nil
		}
		// if it's not a dir or a filepath it's a filename & default to aserto dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to determine user home directory")
	}

	return filepath.Join(home, ".config", x.AppName, ConfigPath), err
}

func setDefaultServices(v *viper.Viper, svc *x.Services) {
	setServiceDefaults(v, "authorizer", &svc.AuthorizerService)
	setServiceDefaults(v, "decision_logs", &svc.DecisionLogsService)
	setServiceDefaults(v, "tenant", &svc.TenantService)
	setServiceDefaults(v, "control_plane", &svc.ControlPlaneService)
	setServiceDefaults(v, "ems", &svc.EMSService)
	setServiceDefaults(v, "directory_reader", &svc.DirectoryReaderService)
	setServiceDefaults(v, "directory_writer", &svc.DirectoryWriterService)
	setServiceDefaults(v, "directory_model", &svc.DirectoryModelService)
}

func setServiceDefaults(v *viper.Viper, name string, svc *x.ServiceOptions) {
	v.SetDefault(fmt.Sprintf("services.%s.address", name), svc.Address)
	v.SetDefault(fmt.Sprintf("services.%s.api_key", name), svc.APIKey)
	v.SetDefault(fmt.Sprintf("services.%s.anonymous", name), svc.Anonymous)
	v.SetDefault(fmt.Sprintf("services.%s.insecure", name), svc.Insecure)
}
