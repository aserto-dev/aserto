package directory

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/pkg/errors"
)

type DeleteUsersCmd struct{}

func (cmd *DeleteUsersCmd) Run(c *cc.CommonCtx) error {
	return errors.Errorf("NOT IMPLEMENTED")

	// client, err := c.AuthorizerClient()
	// if err != nil {
	// 	return err
	// }

	// token := ""
	// pageSize := int32(100)
	// apiUsers := []*api.User{}

	// for {
	// 	resp, listUsersErr := client.Directory.ListUsers(c.Context, &dir.ListUsersRequest{
	// 		Page: &api.PaginationRequest{
	// 			Size:  pageSize,
	// 			Token: token,
	// 		},
	// 		Fields: &api.Fields{
	// 			Mask: &fieldmaskpb.FieldMask{
	// 				Paths: []string{"id", "email"},
	// 			},
	// 		},
	// 	})
	// 	if listUsersErr != nil {
	// 		return errors.Wrapf(err, "list users")
	// 	}

	// 	apiUsers = append(apiUsers, resp.Results...)

	// 	if resp.Page.NextToken == "" {
	// 		break
	// 	}

	// 	token = resp.Page.NextToken
	// }

	// fmt.Fprintf(c.UI.Output(), "tenant %s\n", c.TenantID())
	// fmt.Fprintf(c.UI.Output(), "!!! deleting %d users\n", len(apiUsers))
	// fmt.Fprintf(c.UI.Output(), "please acknowledge that is what you want by typing \"CONFIRMED\" (all uppercase)\n")
	// var input string
	// n, err := fmt.Fscanln(os.Stdin, &input)
	// if err != nil || n == 0 {
	// 	return err
	// }
	// if input == "CONFIRMED" {
	// 	fmt.Fprintf(c.UI.Output(), "starting deletion\n")
	// 	for i, u := range apiUsers {
	// 		fmt.Fprintf(os.Stderr, "\033[2K\rdeleted %d of %d", i+1, len(apiUsers))
	// 		if _, err := client.Directory.DeleteUser(c.Context, &dir.DeleteUserRequest{
	// 			Id: u.Id,
	// 		}); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	fmt.Fprintln(os.Stderr)
	// }

	// return nil
}
