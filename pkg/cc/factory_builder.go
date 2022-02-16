package cc

import (
	"log"

	"github.com/aserto-dev/aserto/pkg/client"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

var (
	NeedLoginErr    = errors.Errorf("user is not logged in, please login using '%s login'", x.AppName)
	TokenExpiredErr = errors.Errorf("the access token has expired, please login using '%s login'", x.AppName)
	NeedTenantIDErr = errors.Errorf("operation requires tenant-id, please login using '%s login' or use --tenant to specify an id.", x.AppName)
)

type ClientFactoryBuilder struct {
	needsToken bool
	overrides  map[x.Service]client.Overrides
}

func NewClientFactoryBuilder() *ClientFactoryBuilder {
	return &ClientFactoryBuilder{overrides: map[x.Service]client.Overrides{}}
}

func (b *ClientFactoryBuilder) Override(svc x.Service, overrides client.Overrides) {
	b.overrides[svc] = overrides
}

func (b *ClientFactoryBuilder) RequireToken() {
	b.needsToken = true
}

func (b *ClientFactoryBuilder) ClientFactory(env *x.Services, tenantID string, token CachedToken) (*client.AsertoFactory, error) {
	if b.needsToken {
		log.Print("command requires token. verifying...")
		if err := token.Verify(); err != nil {
			return nil, err
		}
	}

	opts := map[x.Service]client.OptionsBuilder{}
	for _, svc := range x.AllServices {
		config := &optionsBuilder{
			Overrides:   b.serviceOverrides(svc),
			service:     svc,
			defaultAddr: env.AddressOf(svc),
			tenantID:    tenantID,
			token:       token,
		}

		opts[svc] = config.ConnectionOptions
	}

	return &client.AsertoFactory{SvcOptions: opts}, nil
}

func (b *ClientFactoryBuilder) serviceOverrides(svc x.Service) client.Overrides {
	overrides, ok := b.overrides[svc]
	if !ok {
		return &noOverrides{}
	}

	return overrides
}

type noOverrides struct{}

func (o *noOverrides) Address() string {
	return ""
}

func (o *noOverrides) Key() string {
	return ""
}

func (o *noOverrides) IsAnonymous() bool {
	return false
}

func (o *noOverrides) IsInsecure() bool {
	return false
}
