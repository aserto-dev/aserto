package client

// Overrides are options that modify the default behavior of connections to aserto services.
type Overrides interface {
	// Address overrides the default address used to connect to a service.
	Address() string

	// Key provides an API key to be used instead of the default access token.
	Key() string

	// IsAnonymous means no credentials are sent to the service.
	IsAnonymous() bool

	// IsInsecure indicates that no TLS verification is performed.
	IsInsecure() bool
}
