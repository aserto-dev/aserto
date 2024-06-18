package errors

import (
	"github.com/aserto-dev/aserto/pkg/x"
	"github.com/pkg/errors"
)

var (
	NeedLoginErr     = errors.Errorf("user is not logged in, please login using '%s login'", x.AppName)
	TokenExpiredErr  = errors.Errorf("the access token has expired, please login using '%s login'", x.AppName)
	NeedTenantIDErr  = errors.Errorf("operation requires tenant-id, please login using '%s login' or use --tenant to specify an id.", x.AppName)
	ResolveTenantErr = errors.New("cannot resolve tenant name %q to tenant ID")
)
