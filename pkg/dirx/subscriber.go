package dirx

import (
	"context"

	"github.com/aserto-dev/proto/aserto/api"
	dir "github.com/aserto-dev/proto/aserto/authorizer/directory"
	"github.com/pkg/errors"
)

type Result struct {
	Counts *dir.LoadUsersResponse
	Err    error
}

// Subscriber subscribes to the api.User channel and sends the users instance to the directory using the gRPC LoadUsers API.
func Subscriber(ctx context.Context, client dir.DirectoryClient, s <-chan *api.User, r chan<- *Result, errc chan<- error, inclAttrSets bool) {

	stream, err := client.LoadUsers(ctx)
	if err != nil {
		errc <- errors.Wrapf(err, "client.LoadUsers")
	}

	sendCount := int32(0)
	errCount := int32(0)

	for user := range s {
		if !inclAttrSets {
			user.Attributes = &api.AttrSet{}
			user.Applications = make(map[string]*api.AttrSet)
		}

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

	r <- &Result{
		Counts: res,
		Err:    err,
	}
}

// UserExtSubscriber subscribes to the api.User channel and sends user extensions (api.UserExt message) to the directory using the gRPC LoadUser API.
func UserExtSubscriber(ctx context.Context, client dir.DirectoryClient, s <-chan *api.User, r chan<- *Result, errc chan<- error) {

	stream, err := client.LoadUsers(ctx)
	if err != nil {
		errc <- errors.Wrapf(err, "client.LoadUsers")
	}

	sendCount := int32(0)
	errCount := int32(0)

	for user := range s {

		userExt := api.UserExt{
			Id:           user.Id,
			Attributes:   user.Attributes,
			Applications: user.Applications,
		}

		req := &dir.LoadUsersRequest{
			Data: &dir.LoadUsersRequest_UserExt{
				UserExt: &userExt,
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

	r <- &Result{
		Counts: res,
		Err:    err,
	}
}
