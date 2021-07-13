# Aserto CLI Todo List

## Todo

- [ ] Add localhost flag for command which can target a local onebox, like directory load-users --local 
- [ ] Validate (account) ids using go-lib/ids validation mechanism
- [ ] Stop using "latest" container image for dev mode; use specific tag and semver mask which ideally matches the hosted authorizer or runtime version to guarantee compatibility of local onebox with hosted authorizer instance (would require discovery of authorizer runtime and proto versions)
- [ ] Remove database from dev install, use users.json file instead as it has multiple usage scenarios and is less susceptible  for schema changes
- [ ] Add launch hosted console command
- [ ] Change keyring dependency to on which does not require CGO, like [github.com/docker/docker-credential-helpers](github.com/docker/docker-credential-helpers)
- [ ] Implement device code flow authentication for headless usage
- [ ] Authorizer exec-query add flags for metrics and trace
- [ ] Authorizer exec-query add ability to provide --input from file beside string
## Done

- [X] Store UTC token expiration timestamp (creation timestamp + expires_in) to determine when token is expire before making a call, so that we can provide a better user experience (not having to redo a call because they need to login after token expired failure)
- [X] Config set-env: making env persistent (prod, eng)
- [X] Change get|set|del-user-ext to use positional arguments instead of --id for identifying the user id.
- [X] Change get|set|del-appl-ext to use positional arguments instead of --id and --name for identifying the user id and application name. 
- [X] Config set-tenant: allow user to set tenant context within the collection of tenants of the users (tenant.GetAccount)
- [X] Directory load-users: fix overloaded use of --input, split into --profile and --file (for data input files)
- [X] Change exec-query to accept --stdin | --file as query input
- [X] Make keychain environment independent
- [X] Change dir get-user to be positional argument based instead of --id
- [X] Change dir get-identity to be positional argument instead of --identity
- [X] Change dir get-user to accepts user ids and user identities as inputs <user id | identity>
- [X] Change config get to be positional for property names
- [X] Fix regression introduced in v0.0.17 tenant and authorizer override 
- [X] Rename tenant list-policy-packages to list-policy-references
