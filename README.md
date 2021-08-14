# aserto

Aserto CLI

```
Usage: aserto <command>

Welcome to modern authorization!

Commands:
  login      login
  logout     logout
  version    version information

authorizer
  authorizer eval-decision    evaluate policy decision
  authorizer decision-tree    get decision tree
  authorizer exec-query       execute query

tenant
  tenant get-account               get account info
  tenant list-connections          list connections
  tenant get-connection            get connection instance info
  tenant verify-connection         verify connection settings
  tenant sync-connection           trigger sync of IDP connection
  tenant list-policy-references    list policy references
  tenant create-policy-push-key    create policy upload key
  tenant list-provider-kinds       list provider kinds
  tenant list-providers            list providers
  tenant get-provider              get provider info

identity
  directory get-identity     resolve user identity
  directory list-users       list users
  directory get-user         retrieve user object
  directory load-users       load users
  directory load-user-ext    load user extensions
  directory set-user         disable|enable user
  directory delete-users     delete users from edge directory

user extensions
  directory get-user-props    get properties
  directory set-user-prop     set property
  directory del-user-prop     delete property
  directory get-user-roles    get roles
  directory set-user-role     set role
  directory del-user-role     delete role
  directory get-user-perms    get permissions
  directory set-user-perm     set permission
  directory del-user-perm     delete permission

user application extensions
  directory list-user-apps    list user applications
  directory set-user-app      set user application
  directory del-user-app      delete user application
  directory get-appl-props    get properties
  directory set-appl-prop     set property
  directory del-appl-prop     delete property
  directory get-appl-roles    get roles
  directory set-appl-role     set role
  directory del-appl-role     delete role
  directory get-appl-perms    get permissions
  directory set-appl-perm     set permission
  directory del-appl-perm     delete permission

tenant resources
  directory list-res    list resources
  directory get-res     get resource
  directory set-res     set resource
  directory del-res     delete resource

developer
  developer start        start aserto-one instance
  developer stop         stop aserto-one instance
  developer status       status of aserto-one instance
  developer update       download the latest aserto onebox image
  developer console      launch web console
  developer configure    configure a policy
  developer install      install aserto onebox
  developer uninstall    uninstall aserto onebox, removes all locally installed artifacts

user
  user info    get user profile information
  user get     get property

config
  config get-tenant    get tenant list
  config set-tenant    set default tenant
  config get-env       get environment info

Flags:
  -h, --help                 Show context-sensitive help.
      --verbose              verbose output
      --authorizer=STRING    authorizer override ($ASERTO_AUTHORIZER)
      --tenant=STRING        tenant id override ($ASERTO_TENANT_ID)
      --debug                enable debug logging ($ASERTO_DEBUG)

Run "aserto <command> --help" for more information on a command.
```