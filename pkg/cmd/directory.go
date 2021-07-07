package cmd

import (
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/handlers/directory"
)

type DirectoryCmd struct {
	GetIdentity directory.GetIdentityCmd `cmd:"" help:"resolve user identity" group:"identity"`
	ListUsers   directory.ListUsersCmd   `cmd:"" help:"list users" group:"identity"`
	GetUser     directory.GetUserCmd     `cmd:"" help:"retrieve user object" group:"identity"`
	LoadUsers   directory.LoadUsersCmd   `cmd:"" help:"load users" group:"identity"`
	SetUser     directory.SetUserCmd     `cmd:"" help:"disable|enable user" group:"identity"`
	DeleteUsers directory.DeleteUsersCmd `cmd:"" help:"delete users from edge directory" group:"identity"`

	GetUserProps directory.GetUserPropsCmd `cmd:"" help:"get properties" group:"user extensions"`
	SetUserProp  directory.SetUserPropCmd  `cmd:"" help:"set property" group:"user extensions"`
	DelUserProp  directory.DelUserPropCmd  `cmd:"" help:"delete property" group:"user extensions"`
	GetUserRoles directory.GetUserRolesCmd `cmd:"" help:"get roles" group:"user extensions"`
	SetUserRole  directory.SetUserRoleCmd  `cmd:"" help:"set role" group:"user extensions"`
	DelUserRole  directory.DelUserRoleCmd  `cmd:"" help:"delete role" group:"user extensions"`
	GetUserPerms directory.GetUserPermsCmd `cmd:"" help:"get permissions" group:"user extensions"`
	SetUserPerm  directory.SetUserPermCmd  `cmd:"" help:"set permission" group:"user extensions"`
	DelUserPerm  directory.DelUserPermCmd  `cmd:"" help:"delete permission" group:"user extensions"`

	ListUserApps directory.ListUserAppsCmd `cmd:"" help:"list user applications" group:"user application extensions"`
	SetUserApp   directory.SetUserAppCmd   `cmd:"" help:"set user application" group:"user application extensions"`
	DelUserApp   directory.DelUserAppCmd   `cmd:"" help:"delete user application" group:"user application extensions"`
	GetApplProps directory.GetApplPropsCmd `cmd:"" help:"get properties" group:"user application extensions"`
	SetApplProp  directory.SetApplPropCmd  `cmd:"" help:"set property" group:"user application extensions"`
	DelApplProp  directory.DelApplPropCmd  `cmd:"" help:"delete property" group:"user application extensions"`
	GetApplRoles directory.GetApplRolesCmd `cmd:"" help:"get roles" group:"user application extensions"`
	SetApplRole  directory.SetApplRoleCmd  `cmd:"" help:"set role" group:"user application extensions"`
	DelApplRole  directory.DelApplRoleCmd  `cmd:"" help:"delete role" group:"user application extensions"`
	GetApplPerms directory.GetApplPermsCmd `cmd:"" help:"get permissions" group:"user application extensions"`
	SetApplPerm  directory.SetApplPermCmd  `cmd:"" help:"set permission" group:"user application extensions"`
	DelApplPerm  directory.DelApplPermCmd  `cmd:"" help:"delete permission" group:"user application extensions"`

	ListRes directory.ListResCmd `cmd:"" help:"list resources" group:"tenant resources"`
	GetRes  directory.GetResCmd  `cmd:"" help:"get resource" group:"tenant resources"`
	SetRes  directory.SetResCmd  `cmd:"" help:"set resource" group:"tenant resources"`
	DelRes  directory.DelResCmd  `cmd:"" help:"delete resource" group:"tenant resources"`
}

func (cmd *DirectoryCmd) BeforeApply(c *CLI) error {
	c.RequireLogin()
	return nil
}

func (cmd *DirectoryCmd) Run(c *cc.CommonCtx) error {
	return nil
}
