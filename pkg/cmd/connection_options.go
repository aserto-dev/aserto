package cmd

type ConnectionOptions struct {
	APIKey     string `env:"KEY" help:"service api key" placeholder:"key"`
	NoAuth     bool   `help:"do not provide any credentials"`
	Insecure   bool   `help:"skip TLS verification" default:"false"`
	CACertPath string `help:"path to grpc CA cert"`
}

func (so *ConnectionOptions) Address() string {
	return ""
}

func (so *ConnectionOptions) Key() string {
	return so.APIKey
}

func (so *ConnectionOptions) IsAnonymous() bool {
	return so.NoAuth
}

func (so *ConnectionOptions) IsInsecure() bool {
	return so.Insecure
}

func (so *ConnectionOptions) PathToCACert() string {
	return so.CACertPath
}

type ServiceOverrideOptions struct {
	AddressOverride string `name:"address" env:"ADDRESS" help:"address override" default:""`

	ConnectionOptions
}

func (ao *ServiceOverrideOptions) Address() string {
	return ao.AddressOverride
}
