package x

import (
	"log"

	"github.com/pkg/errors"
)

type Services struct {
	Environment         string
	TenantService       string
	AuthorizerService   string
	DecisionLogsService string
}

type Service int

const (
	AuthorizerService Service = iota
	DecisionLogsService
	TenantService
)

var (
	serviceNames = map[Service]string{
		AuthorizerService:   "authorizer",
		DecisionLogsService: "decision logs",
		TenantService:       "tenant",
	}

	AllServices = []Service{AuthorizerService, DecisionLogsService, TenantService}
)

func (s Service) Name() string {
	name, ok := serviceNames[s]
	if ok {
		return name
	}

	return "unknown"
}

func (s *Services) AddressOf(svc Service) string {
	switch svc {
	case AuthorizerService:
		return s.AuthorizerService
	case DecisionLogsService:
		return s.DecisionLogsService
	case TenantService:
		return s.TenantService
	}

	log.Panicf("unknown service [%d]\n", svc)
	return ""
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
