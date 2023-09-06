package auth0

import (
	"strings"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/auth0/go-auth0/management"
)

const (
	provider = "auth0"
)

// Transform Auth0 user definition into Aserto Edge User object definition.
func Transform(in *management.User) (*api.User, error) {

	uid := strings.ToLower(strings.TrimPrefix(*in.ID, "auth0|"))

	user := api.User{
		Id:          uid,
		DisplayName: in.GetNickname(),
		Email:       in.GetEmail(),
		Picture:     in.GetPicture(),
		Identities:  make(map[string]*api.IdentitySource),
		Attributes: &api.AttrSet{
			Properties:  &structpb.Struct{Fields: make(map[string]*structpb.Value)},
			Roles:       []string{},
			Permissions: []string{},
		},
		Applications: make(map[string]*api.AttrSet),
		Metadata: &api.Metadata{
			CreatedAt: timestamppb.New(in.GetCreatedAt()),
			UpdatedAt: timestamppb.New(in.GetUpdatedAt()),
		},
	}

	user.Identities[in.GetID()] = &api.IdentitySource{
		Kind:     api.IdentityKind_IDENTITY_KIND_PID,
		Provider: provider,
		Verified: true,
	}

	user.Identities[in.GetEmail()] = &api.IdentitySource{
		Kind:     api.IdentityKind_IDENTITY_KIND_EMAIL,
		Provider: provider,
		Verified: in.GetEmailVerified(),
	}

	if in.UserMetadata != nil && len(*in.UserMetadata) != 0 {
		props, err := structpb.NewStruct(*in.UserMetadata)
		if err == nil {
			user.Attributes.Properties = props
		}
	}

	return &user, nil
}
