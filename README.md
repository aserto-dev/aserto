# aserto

Aserto CLI

```
Usage: aserto <command>

Welcome to modern authorization!

Commands:
  developer (xp)        developer commands
  authorizer (az)       authorizer commands
  decision-logs (dl)    decision logs commands
  control-plane (cp)    control plane commands
  tenant (tn)           tenant commands
  login                 login
  logout                logout
  config                configuration commands
  version               version information

Flags:
  -h, --help             Show context-sensitive help.
  -v, --verbosity=INT    Use to increase output verbosity.

Run "aserto <command> --help" for more information on a command.
```

In order to use the commands you have to call `aserto login` first. The `developer` and
`authorizer` commands can operate without a prior log in.

## Configuration commands

After log in, you can view your tenants by calling `aserto config get-tenants` and you can switch between your tenants using the `use-tenant` command.

Contexts refer to an authorizer env. If you start a sidecar using the `developer`
commands, you can run the `authorizer` commands agains it by defining a new context
using `set-context` and switching to it by `use-context` command. 

```
config
  config user-info              get user profile information
  config get-property           get property
  config create-tenant-alias    create an alias for a tenant name
  config get-tenants            get defined tenants
  config use-tenant             use a specific tenant
  config refresh-tenant         refetch api keys for a specific tenant context
  config current-tenant         get the tenant name in use
  config get-contexts           get defined contexts
  config get-active-context     get active context config
  config delete-context         delete a context config
  config set-context            creates a context
  config use-context            use a specific context

Flags:
  -h, --help             Show context-sensitive help.
  -v, --verbosity=INT    Use to increase output verbosity.
```

## Developer commands 

```
developer
  developer (xp) start                   start sidecar instance
  developer (xp) stop                    stop sidecar instance
  developer (xp) status                  status of sidecar instance
  developer (xp) update                  download the latest aserto sidecar image
  developer (xp) console                 launch web console
  developer (xp) configure               configure a policy
  developer (xp) policy-from-open-api    generate an open api policy
  developer (xp) install                 install aserto sidecar
  developer (xp) uninstall               uninstall aserto sidecar, removes all locally installed artifacts

Flags:
  -h, --help             Show context-sensitive help.
  -v, --verbosity=INT    Use to increase output verbosity.
```

## Authorizer commands

```
authorizer
  authorizer (az) eval-decision    evaluate policy decision
  authorizer (az) decision-tree    get decision tree
  authorizer (az) exec-query       execute query
  authorizer (az) compile          compile query
  authorizer (az) get-policy       get policy
  authorizer (az) list-policies    list policies

Flags:
  -h, --help             Show context-sensitive help.
  -v, --verbosity=INT    Use to increase output verbosity.

      --authorizer=""    authorizer override ($ASERTO_AUTHORIZER_ADDRESS)
      --api-key=key      service api key ($ASERTO_AUTHORIZER_KEY)
      --no-auth          do not provide any credentials
      --insecure         skip TLS verification
```

## Decision logs commands

```
decision-logs
  decision-logs (dl) list          list available decision log files
  decision-logs (dl) get           download one or more decision log files
  decision-logs (dl) list-users    list available user data files
  decision-logs (dl) get-user      download one or more user data files
  decision-logs (dl) stream        stream decision log events to stdout

Flags:
  -h, --help             Show context-sensitive help.
  -v, --verbosity=INT    Use to increase output verbosity.

      --api-key=key      service api key ($ASERTO_DECISION_LOGS_KEY)
      --no-auth          do not provide any credentials
      --insecure         skip TLS verification
```

## Control plane commands

```
control-plane
  control-plane (cp) list-connections    list edge authorizer connections
  control-plane (cp) client-cert         get client certificates for an edge authorizer connection
  control-plane (cp) list-instance-registrations
                                         list instance registrations
  control-plane (cp) discovery           run discovery on a registered instance
  control-plane (cp) edge-dir-sync       sync the directory on an edge authorizer

Flags:
  -h, --help             Show context-sensitive help.
  -v, --verbosity=INT    Use to increase output verbosity.
```

## Tenant commands

```
tenant
  tenant (tn) get-account               get account info
  tenant (tn) list-connections          list connections
  tenant (tn) get-connection            get connection instance info
  tenant (tn) update-connection         update connection configuration fields
  tenant (tn) verify-connection         verify connection settings
  tenant (tn) sync-connection           trigger sync of IDP connection
  tenant (tn) list-policy-references    list policy references
  tenant (tn) list-provider-kinds       list provider kinds
  tenant (tn) list-providers            list providers
  tenant (tn) get-provider              get provider info

Flags:
  -h, --help             Show context-sensitive help.
  -v, --verbosity=INT    Use to increase output verbosity.
```