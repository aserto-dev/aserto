package cmd

type ConnectionOptions struct {
	APIKey   string `env:"KEY" help:"service api key" placeholder:"key"`
	NoAuth   bool   `help:"do not provide any credentials"`
	Insecure bool   `help:"skip TLS verification" default:"false"`
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

type AuthorizerOptions struct {
	AddressOverride string `name:"authorizer" env:"ADDRESS" help:"authorizer override" default:""`

	ConnectionOptions
}
