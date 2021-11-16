package json

import (
	"encoding/json"
	"os"

	"github.com/aserto-dev/aserto/pkg/pb"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"github.com/pkg/errors"
)

// Producer api.User producer.
type Producer struct {
	filename string
	count    int
}

func (p *Producer) Count() int {
	return p.count
}

// NewProducer returns Auth0 producer instance.
func NewProducer(filename string) *Producer {
	return &Producer{
		filename: filename,
	}
}

// Producer creates a stream of api.User instances from a JSON input file.
func (p *Producer) Producer(s chan *api.User, errc chan error) {

	r, err := os.Open(p.filename)
	if err != nil {
		errc <- errors.Wrapf(err, "open %s", p.filename)
	}

	dec := json.NewDecoder(r)

	if _, err = dec.Token(); err != nil {
		errc <- errors.Wrapf(err, "token open")
	}

	for dec.More() {
		u := api.User{}
		if err = pb.UnmarshalNext(dec, &u); err != nil {
			errc <- errors.Wrapf(err, "unmarshal next")
		}

		s <- &u
		p.count++
	}

	if _, err = dec.Token(); err != nil {
		errc <- errors.Wrapf(err, "token close")
	}
}
