package dirx

import (
	"context"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"
	"github.com/pkg/errors"
)

type Result struct {
	Counts *dir.LoadUsersResponse
	Err    error
}

type UserSubscriber struct {
	Ctx           context.Context
	DirClient     dir.DirectoryClient
	SourceChannel chan *api.User
	ResultChannel chan *Result
	ErrorChannel  chan error
}

type LoadUsersRequestFactory func(*api.User) *dir.LoadUsersRequest

func (subscriber *UserSubscriber) Subscribe(requestFactory LoadUsersRequestFactory) {
	stream, err := subscriber.DirClient.LoadUsers(subscriber.Ctx)
	if err != nil {
		subscriber.ErrorChannel <- errors.Wrapf(err, "client.LoadUsers")
	}

	sendCount := int32(0)
	errCount := int32(0)

	for user := range subscriber.SourceChannel {
		loadUsersRequest := requestFactory(user)

		if err = stream.Send(loadUsersRequest); err != nil {
			subscriber.ErrorChannel <- errors.Wrapf(err, "stream send %s", user.Id)
			errCount++
		}
		sendCount++
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		subscriber.ErrorChannel <- errors.Wrapf(err, "stream.CloseAndRecv()")
	}

	subscriber.ResultChannel <- &Result{
		Counts: res,
		Err:    err,
	}
}

func NewLoadUsersRequestFactory(inclAttrSets bool) LoadUsersRequestFactory {
	return func(user *api.User) *dir.LoadUsersRequest {
		if !inclAttrSets {
			user.Attributes = &api.AttrSet{}
			user.Applications = make(map[string]*api.AttrSet)
		}

		return &dir.LoadUsersRequest{
			Data: &dir.LoadUsersRequest_User{
				User: user,
			},
		}
	}
}

func UserExtensionsRequestFactory(user *api.User) *dir.LoadUsersRequest {
	userExt := api.UserExt{
		Id:           user.Id,
		Attributes:   user.Attributes,
		Applications: user.Applications,
	}

	return &dir.LoadUsersRequest{
		Data: &dir.LoadUsersRequest_UserExt{
			UserExt: &userExt,
		},
	}
}
