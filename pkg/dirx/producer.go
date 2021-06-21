package dirx

import "github.com/aserto-dev/proto/aserto/api"

// Producer interface.
type Producer interface {
	Producer(chan<- *api.User, chan<- error)
	Count() int
}
