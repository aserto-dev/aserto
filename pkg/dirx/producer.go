package dirx

import api "github.com/aserto-dev/go-grpc/aserto/api/v1"

// Producer interface.
type Producer interface {
	Producer(chan<- *api.User, chan<- error)
	Count() int
}
