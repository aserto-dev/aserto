package authorizer

import (
	"encoding/json"
	"fmt"

	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

type IdentityType string

const (
	IdentityTypeNone IdentityType = "none"
	IdentityTypeSub  IdentityType = "sub"
	IdentityTypeJwt  IdentityType = "jwt"
)

type AuthParams struct {
	Identity     string       `name:"identity" help:"caller identity" default:""`
	IdentityType IdentityType `name:"identity-type" enum:"sub,jwt,none" help:"type of identity [sub|jwt|none]"  default:"none"`
	Resource     string       `name:"resource" help:"a JSON object to include as resource context"`
	PolicyID     string       `name:"policy-id" required:"" help:"policy id"`
}

func (a AuthParams) IdentityContext() *api.IdentityContext {
	id_type := api.IdentityType_IDENTITY_TYPE_NONE
	switch a.IdentityType {
	case IdentityTypeSub:
		id_type = api.IdentityType_IDENTITY_TYPE_SUB
	case IdentityTypeJwt:
		id_type = api.IdentityType_IDENTITY_TYPE_JWT
	}

	return &api.IdentityContext{
		Identity: a.Identity,
		Type:     id_type,
	}
}

func (a AuthParams) ResourceContext() (*structpb.Struct, error) {
	result := &structpb.Struct{}

	if a.Resource != "" {
		var r interface{}
		if err := json.Unmarshal([]byte(a.Resource), &r); err != nil {
			return result, err
		}

		m, ok := r.(map[string]interface{})
		if !ok {
			return result, fmt.Errorf("resource must be a JSON object")
		}

		return structpb.NewStruct(m)
	}

	return result, nil
}
