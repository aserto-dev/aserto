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
	EMSService
)

var (
	UnknownSvcErr = errors.New("unknown service")

	serviceNames = map[Service]string{
		AuthorizerService:   "authorizer",
		DecisionLogsService: "decision logs",
		TenantService:       "tenant",
		ControlPlaneService: "control plane",
		EMSService:          "ems",
	}

	AllServices = []Service{AuthorizerService, DecisionLogsService, TenantService, ControlPlaneService, EMSService}
)

func (s Service) Name() string {
	name, ok := serviceNames[s]
	if ok {
		return name
	}

	return "unknown"
}

type ServiceOptions struct {
	Address   string `json:"address,omitempty" yaml:"address,omitempty"`
	APIKey    string `json:"api_key,omitempty" yaml:"api_key,omitempty"`
	Anonymous bool   `json:"anonymous,omitempty" yaml:"anonymous,omitempty"`
	Insecure  bool   `json:"insecure,omitempty" yaml:"insecure,omitempty"`
}

type Services struct {
	DecisionLogsService ServiceOptions `json:"decision_logs" yaml:"decision_logs"`
	TenantService       ServiceOptions `json:"tenant" yaml:"tenant"`
	ControlPlaneService ServiceOptions `json:"control_plane" yaml:"control_plane"`
	EMSService          ServiceOptions `json:"ems" yaml:"ems"`
	AuthorizerService   ServiceOptions `json:"authorizer" yaml:"authorizer"`
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
	case EMSService:
		return &s.EMSService
	default:
		log.Panicf("unknown service [%d]\n", svc)
	}

	return &ServiceOptions{}
}

func (s *Services) SetAddress(svc Service, address string) error {
	switch svc { //nolint:exhaustive
	case AuthorizerService:
		s.AuthorizerService.Address = address
	case DecisionLogsService:
		s.DecisionLogsService.Address = address
	case TenantService:
		s.TenantService.Address = address
	case ControlPlaneService:
		s.ControlPlaneService.Address = address
	case EMSService:
		s.EMSService.Address = address
	default:
		return errors.Wrapf(UnknownSvcErr, "[%d]", svc)
	}

	return nil
}

func DefaultEnvironment() *Services {
	return &Services{
		AuthorizerService:   ServiceOptions{Address: "authorizer.prod.aserto.com:8443"},
		TenantService:       ServiceOptions{Address: "tenant.prod.aserto.com:8443"},
		DecisionLogsService: ServiceOptions{Address: "decision-logs.prod.aserto.com:8443"},
		ControlPlaneService: ServiceOptions{Address: "relay.prod.aserto.com:8443"},
		EMSService:          ServiceOptions{Address: "ems.prod.aserto.com:8443"},
	}
}
