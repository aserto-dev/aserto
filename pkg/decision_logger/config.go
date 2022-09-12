package decisionlogger

const (
	Dir           = "decision_logs"
	ContainerPath = "/app/" + Dir
)

type Config struct {
	EMSAddress string `json:"ems_address"`
}

type Settings Config

func NewSettings(cfg *Config) *Settings {
	s := Settings(*cfg)
	if s.EMSAddress == "" {
		s.EMSAddress = "ems.prod.aserto.com:8443"
	}

	return &s
}
