package x

import (
	"github.com/pkg/errors"
)

type Services struct {
	Environment       string
	TenantService     string
	AuthorizerService string
	RegistryService   string
	TasksService      string
}

func Environment(env string) (*Services, error) {
	switch env {
	case EnvProduction:
		return &Services{
			Environment:       EnvProduction,
			TenantService:     "tenant.prod.aserto.com:8443",
			AuthorizerService: "authorizer.prod.aserto.com:8443",
			RegistryService:   "https://bundler.prod.aserto.com",
		}, nil
	case EnvEngineering:
		return &Services{
			Environment:       EnvEngineering,
			TenantService:     "tenant.eng.aserto.com:8443",
			AuthorizerService: "authorizer.eng.aserto.com:8443",
			RegistryService:   "https://bundler.eng.aserto.com",
		}, nil
	default:
		return nil, errors.Errorf("invalid environment [%s]", env)
	}
}
