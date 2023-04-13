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
  authorizer (a) eval-decision    evaluate policy decision
  authorizer (a) decision-tree    get decision tree
  authorizer (a) exec-query       execute query

tenant
  tenant (t) get-account          get account info
  tenant (t) list-connections     list connections
  tenant (t) get-connection       get connection instance info
  tenant (t) update-connection    update connection configuration fields
  tenant (t) verify-connection    verify connection settings
  tenant (t) sync-connection      trigger sync of IDP connection
  tenant (t) list-policy-references
                                  list policy references
  tenant (t) list-provider-kinds
                                  list provider kinds
  tenant (t) list-providers       list providers
  tenant (t) get-provider         get provider info

identity
  directory (d) get-identity     resolve user identity
  directory (d) list-users       list users
  directory (d) get-user         retrieve user object
  directory (d) load-users       load users
  directory (d) load-user-ext    load user extensions
  directory (d) set-user         disable|enable user
  directory (d) delete-users     delete users from edge directory

user extensions
  directory (d) get-user-props    get properties
  directory (d) set-user-prop     set property
  directory (d) del-user-prop     delete property
  directory (d) get-user-roles    get roles
  directory (d) set-user-role     set role
  directory (d) del-user-role     delete role
  directory (d) get-user-perms    get permissions
  directory (d) set-user-perm     set permission
  directory (d) del-user-perm     delete permission

user application extensions
  directory (d) list-user-apps    list user applications
  directory (d) set-user-app      set user application
  directory (d) del-user-app      delete user application
  directory (d) get-appl-props    get properties
  directory (d) set-appl-prop     set property
  directory (d) del-appl-prop     delete property
  directory (d) get-appl-roles    get roles
  directory (d) set-appl-role     set role
  directory (d) del-appl-role     delete role
  directory (d) get-appl-perms    get permissions
  directory (d) set-appl-perm     set permission
  directory (d) del-appl-perm     delete permission

tenant resources
  directory (d) list-res    list resources
  directory (d) get-res     get resource
  directory (d) set-res     set resource
  directory (d) del-res     delete resource

decision-logs
  decision-logs (l) list          list available decision log files
  decision-logs (l) get           download one or more decision log files
  decision-logs (l) list-users    list available user data files
  decision-logs (l) get-user      download one or more user data files

developer
  developer (x) start        start sidecar instance
  developer (x) stop         stop sidecar instance
  developer (x) status       status of sidecar instance
  developer (x) update       download the latest aserto sidecar image
  developer (x) console      launch web console
  developer (x) configure    configure a policy
  developer (x) install      install aserto sidecar
  developer (x) uninstall    uninstall aserto sidecar, removes all locally
                             installed artifacts

user
  user (u) info    get user profile information
  user (u) get     get property

config
  config (c) get-tenant    get tenant list
  config (c) set-tenant    set default tenant
  config (c) get-env       get environment info

Flags:
  -h, --help             Show context-sensitive help.
  -c, --config=STRING    name or path of configuration file ($ASERTO_ENV)
  -v, --verbosity=INT    Use to increase output verbosity.
      --tenant=STRING    tenant id override ($ASERTO_TENANT_ID)

Run "aserto <command> --help" for more information on a command.
```
