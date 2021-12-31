package x

import (
	"github.com/pkg/errors"
)

type Services struct {
	Environment         string
	TenantService       string
	AuthorizerService   string
	TasksService        string
	DecisionLogsService string
}

func Environment(env string) (*Services, error) {
	switch env {
	case EnvProduction:
		return &Services{
			Environment:         EnvProduction,
			TenantService:       "tenant.prod.aserto.com:8443",
			AuthorizerService:   "authorizer.prod.aserto.com:8443",
			DecisionLogsService: "decision-logs.prod.aserto.com:8443",
		}, nil
	case EnvEngineering:
		return &Services{
			Environment:         EnvEngineering,
			TenantService:       "tenant.eng.aserto.com:8443",
			AuthorizerService:   "authorizer.eng.aserto.com:8443",
			DecisionLogsService: "decision-logs.eng.aserto.com:8443",
		}, nil
	case EnvBeta:
		return &Services{
			Environment:         EnvBeta,
			TenantService:       "tenant.beta.aserto.com:8443",
			AuthorizerService:   "authorizer.beta.aserto.com:8443",
			DecisionLogsService: "decision-logs.beta.aserto.com:8443",
		}, nil
	default:
		return nil, errors.Errorf("invalid environment [%s]", env)
	}
}
