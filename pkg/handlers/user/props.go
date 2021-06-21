package user

import (
	"fmt"
	"strings"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

type GetCmd struct {
	Name string `arg:"name" name:"name" required:"" help:"property name"`
}

func (cmd *GetCmd) Run(c *cc.CommonCtx) error {

	propName := strings.ToLower(cmd.Name)

	var propValue string
	switch propName {
	case x.PropertyAccessToken:
		propValue = c.AccessToken()
	case x.PropertyTenantID:
		propValue = c.TenantID()
	case x.PropertyAuthorizerAPIKey:
		propValue = c.AuthorizerAPIKey()
	case x.PropertyRegistryDownloadKey:
		propValue = c.RegistryDownloadKey()
	case x.PropertyRegistryUploadKey:
		propValue = c.RegistryUploadKey()
	case x.PropertyToken:
		return jsonx.OutputJSON(c.OutWriter, c.Token())
	default:
		return errors.Errorf("unknown property name [%s]", propName)
	}

	fmt.Fprintf(c.OutWriter, "%s\n", propValue)

	return nil
}
