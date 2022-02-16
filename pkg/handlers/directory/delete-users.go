package directory

import (
	"fmt"
	"os"

	"github.com/aserto-dev/aserto/pkg/cc"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	dir "github.com/aserto-dev/go-grpc/aserto/authorizer/directory/v1"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type DeleteUsersCmd struct{}

func (cmd *DeleteUsersCmd) Run(c *cc.CommonCtx) error {
	client, err := c.AuthorizerClient()
	if err != nil {
		return err
	}

	resp, err := client.Directory.ListUsers(c.Context, &dir.ListUsersRequest{
		Fields: &api.Fields{
			Mask: &fieldmaskpb.FieldMask{
				Paths: []string{"id", "email"},
			},
		},
		Page: &api.PaginationRequest{
			Size: -1,
		},
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(c.OutWriter, "tenant %s\n", c.TenantID())
	fmt.Fprintf(c.OutWriter, "!!! deleting %d users\n", resp.Page.TotalSize)
	fmt.Fprintf(c.OutWriter, "please acknowledge that is what you want by typing \"CONFIRMED\" (all uppercase)\n")
	var input string
	n, err := fmt.Fscanln(os.Stdin, &input)
	if err != nil || n == 0 {
		return err
	}
	if input == "CONFIRMED" {
		fmt.Fprintf(c.OutWriter, "starting deletion\n")
		for i, u := range resp.Results {
			fmt.Fprintf(os.Stderr, "\033[2K\rdeleted %d of %d", i+1, resp.Page.TotalSize)
			if _, err := client.Directory.DeleteUser(c.Context, &dir.DeleteUserRequest{
				Id: u.Id,
			}); err != nil {
				return err
			}
		}
		fmt.Fprintln(os.Stderr)
	}

	return nil
}
