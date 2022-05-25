package x

import (
	"log"

	"github.com/pkg/errors"
)

type Service int

const (
	AuthorizerService Service = iota
	DecisionLogsService
	TenantService
	ControlPlaneService
)

var (
	UnknownSvcErr = errors.New("unknown service")

	serviceNames = map[Service]string{
		AuthorizerService:   "authorizer",
		DecisionLogsService: "decision logs",
		TenantService:       "tenant",
		ControlPlaneService: "control plane",
	}

	AllServices = []Service{AuthorizerService, DecisionLogsService, TenantService, ControlPlaneService}
)

func (s Service) Name() string {
	name, ok := serviceNames[s]
	if ok {
		return name
	}

	return "unknown"
}

type ServiceOptions struct {
	Address   string `json:"address"`
	APIKey    string `json:"api_key,omitempty"`
	Anonymous bool   `json:"anonymous,omitempty"`
	Insecure  bool   `json:"insecure,omitempty"`
}

type Services struct {
	AuthorizerService   ServiceOptions `json:"authorizer"`
	DecisionLogsService ServiceOptions `json:"decision_logs"`
	TenantService       ServiceOptions `json:"tenant"`
	ControlPlaneService ServiceOptions `json:"control_plane"`
}

func (s *Services) Get(svc Service) *ServiceOptions {
	switch svc {
	case AuthorizerService:
		return &s.AuthorizerService
	case DecisionLogsService:
		return &s.DecisionLogsService
	case TenantService:
		return &s.TenantService
	case ControlPlaneService:
		return &s.ControlPlaneService
	default:
		log.Panicf("unknown service [%d]\n", svc)
	}

	return &ServiceOptions{}
}

func (s *Services) SetAddress(svc Service, address string) error {
	switch svc {
	case AuthorizerService:
		s.AuthorizerService.Address = address
	case DecisionLogsService:
		s.DecisionLogsService.Address = address
	case TenantService:
		s.TenantService.Address = address
	case ControlPlaneService:
		s.ControlPlaneService.Address = address
	default:
		return errors.Wrapf(UnknownSvcErr, "[%d]", svc)
	}

	return nil
}

func DefaultEnvironment() *Services {
	return &Services{
		TenantService:       ServiceOptions{Address: "tenant.prod.aserto.com:8443"},
		AuthorizerService:   ServiceOptions{Address: "authorizer.prod.aserto.com:8443"},
		DecisionLogsService: ServiceOptions{Address: "decision-logs.prod.aserto.com:8443"},
		ControlPlaneService: ServiceOptions{Address: "relay.prod.aserto.com:8443"},
	}
}
