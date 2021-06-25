package grpcc

import (
	"github.com/aserto-dev/aserto/pkg/x"
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
	case x.EnvProduction:
		return &Services{
			Environment:       x.EnvProduction,
			TenantService:     "tenant.prod.aserto.com:8443",
			AuthorizerService: "authorizer.prod.aserto.com:8443",
			RegistryService:   "bundler.prod.aserto.com:8443",
			TasksService:      "tasks.prod.aserto.com:8433",
		}, nil
	case x.EnvEngineering:
		return &Services{
			Environment:       x.EnvEngineering,
			TenantService:     "tenant.eng.aserto.com:8443",
			AuthorizerService: "authorizer.eng.aserto.com:8443",
			RegistryService:   "bundler.eng.aserto.com:8443",
			TasksService:      "tasks.eng.aserto.com:8433",
		}, nil
	default:
		return nil, errors.Errorf("invalid environment [%s]", env)
	}
}
