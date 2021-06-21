package auth0

import (
	"github.com/aserto-dev/proto/aserto/api"
	"gopkg.in/auth0.v5/management"
)

// Producer api.User producer.
type Producer struct {
	cfg   *Config
	count int
}

// NewProducer returns Auth0 producer instance.
func NewProducer(cfg *Config) *Producer {
	return &Producer{
		cfg: cfg,
	}
}

func (p *Producer) Count() int {
	return p.count
}

// Producer func.
func (p *Producer) Producer(s chan<- *api.User, errc chan<- error) {
	mgnt, err := management.New(
		p.cfg.Domain,
		management.WithClientCredentials(
			p.cfg.ClientID,
			p.cfg.ClientSecret,
		))
	if err != nil {
		errc <- err
		return
	}

	page := 0
	for {
		opts := management.Page(page)
		ul, err := mgnt.User.List(opts)
		if err != nil {
			errc <- err
		}

		for _, u := range ul.Users {
			user, err := Transform(u)
			if err != nil {
				errc <- err
			}

			s <- user
			p.count++
		}

		if ul.Length < ul.Limit {
			break
		}

		page++
	}
}
