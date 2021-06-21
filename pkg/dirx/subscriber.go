package dirx

import (
	"context"
	"fmt"

	"github.com/aserto-dev/proto/aserto/api"
	dir "github.com/aserto-dev/proto/aserto/authorizer/directory"
	"github.com/pkg/errors"
)

// Subscriber subscribes to the api.User channel and sends the users instance to the directory using the gRPC LoadUsers API.
func Subscriber(ctx context.Context, client dir.DirectoryClient, s <-chan *api.User, done chan<- bool, errc chan<- error) {

	stream, err := client.LoadUsers(ctx)
	if err != nil {
		errc <- errors.Wrapf(err, "client.LoadUsers")
	}

	sendCount := int32(0)
	errCount := int32(0)

	for user := range s {
		req := &dir.LoadUsersRequest{
			Data: &dir.LoadUsersRequest_User{
				User: user,
			},
		}

		if err := stream.Send(req); err != nil {
			errc <- errors.Wrapf(err, "stream send %s", user.Id)
			errCount++
		}
		sendCount++
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		errc <- errors.Wrapf(err, "stream.CloseAndRecv()")
	}

	if res != nil && res.Received != sendCount {
		errc <- fmt.Errorf("send != received %d - %d", sendCount, res.Received)
	}

	done <- true
}
