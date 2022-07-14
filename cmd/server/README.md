# Let's Block It server

As an alternative to [the official instance](https://letsblock.it) and
[the render CLI](https://github.com/letsblockit/letsblockit/blob/main/cmd/render/README.md), experienced administrators
can self-host a server instance themselves. This documentation assumes good knowledge of self-hosting, and "easy install"
methods will not be provided or supported by the project.

You will need:

- The server binary, either the container image or a binary built from this repository,
- A PostgreSQL 14 database and a role with sufficient privileges on it,
- An authentication and authorization provider.

## Running the server

### a. Using the container image

The container image for the server is published as `ghcr.io/letsblockit/server:latest` on the
[GitHub container registry](https://github.com/letsblockit/letsblockit/pkgs/container/server).
We recommend you use the `latest` tag and pull it regularly to get filter updates and fixes.

- All the options can be set via environment variables. Run `docker run ghcr.io/letsblockit/server:latest server --help`
  for an exhaustive list, read the rest of this document for the most important ones.
- While the default value of `LETSBLOCKIT_ADDRESS` is shown as localhost (not serving on public interfaces),
  the container image overrides this with `LETSBLOCKIT_ADDRESS=:8765` (all network interfaces) by default,
  to be reachable from your proxy container.

### b. As a systemd unit

Static binaries are not published for now, but they can be built with `go build ./cmd/server`. 
NixOS users can use the `server` flake output, which powers the official instance. 

<details><summary>Click here for an example systemd unit definition passing options as envvars.</summary>

```ini
[Unit]
After=postgresql.service
Description=letsblock.it server

[Service]
Environment="LETSBLOCKIT_AUTH_PROXY_HEADER_NAME=X-Auth-Request-User"
Environment="LETSBLOCKIT_AUTH_METHOD=proxy"
Environment="LETSBLOCKIT_DATABASE_URL=postgresql:///letsblockit"
ExecStart=/location/to/lbi/server
Restart=always
User=letsblockit
WorkingDirectory=/tmp
```
</details>

- All the options can be set via environment variable. Run `server --help` for an exhaustive list, read the rest of
  this document for the most important ones.
- By default, the server listens to localhost only, on the port `8765`, assuming a reverse-proxy will sit on front
  of it. You can adjust `LETSBLOCKIT_ADDRESS`, or create a systemd socket and set `LETSBLOCKIT_USE_SYSTEMD_SOCKET=true`

## PostgreSQL database

Lists and filter instances are stored in a PostgreSQL 14 database. The project is tested against version 14,
and **support for older versions is not guaranteed**. You should provision a dedicated role and database for the server,
with: `CREATE USER letsblockit; CREATE DATABASE letsblockit OWNER letsblockit;`

The `LETSBLOCKIT_DATABASE_URL` must be set with a valid connection string, usually a URI in the
`postgresql://user:password@host/database` form, see
[the PSQL documentation](https://www.postgresql.org/docs/14/libpq-connect.html#LIBPQ-CONNSTRING) for more details.

The server will create tables and run migrations automatically on startup, thanks to the
[golang-migrate](https://github.com/golang-migrate/migrate) project. You can look for the log lines starting with
`migrate[` during startup. **Rollbacks are not supported**, so we recommend you back up your database before upgrading
the server.

## Authentication and authorization

The server does not include user management, because I do not trust myself to write a secure implementation. Instead,
the official instance relies on [Ory Cloud](https://www.ory.sh/cloud/), and an authenticating reverse-proxy can be used
for self-hosted scenarios. Any authenticating reverse-proxy setup should work, assuming:

- The requests to the `/list/` prefix pass through without authentication, to allow for lists to be downloaded
  by the adblocker
- All the other requests are authenticated, and a unique property of the user (username, email, UUID) is passed
  down as an HTTP Header

The simplest setup would use HTTP Basic auth and a static user list, but you can also use most identity providers
(Authelia, Authentik, Keycloak, Oauth2 Proxy...), either as reverse-proxies or with forward authentication.

Some examples are documented below, but you can [open an issue](https://github.com/letsblockit/letsblockit/issues/new)
for configuration assistance on other setups.

### HTTP basic authentication with Nginx

Create a htpasswd file following [the documentation](https://docs.nginx.com/nginx/admin-guide/security-controls/configuring-http-basic-authentication/)
then configure your vhost accordingly:

```
location /list {  #  No auth for list download
    proxy_pass http://letsblockit:8765;
}
location / {  #  Basic auth for the rest
    auth_basic "Restricted";
    auth_basic_user_file /etc/nginx/htpasswd;
    proxy_pass http://letsblockit:8765;
    proxy_set_header X-Forwarded-User $remote_user;
}
```

You can then run the server with the following variables:
  - `LETSBLOCKIT_AUTH_METHOD=proxy`
  - `LETSBLOCKIT_AUTH_PROXY_HEADER_NAME=X-Forwarded-User`

### HTTP basic authentication with Traefik

You should read the [basicauth middleware doc](https://doc.traefik.io/traefik/middlewares/http/basicauth) and make sure
`headerField` is set. Here is an example envvars and labels for a `letsblockit` container, for Traefik listening on
`localhost`, and the `test:test` and `test2:test2` users:

```yaml
environment:
    - LETSBLOCKIT_AUTH_METHOD=proxy
    - LETSBLOCKIT_AUTH_PROXY_HEADER_NAME=X-Forwarded-User
labels:
    # Allow access to list downloads without authentication
    - "traefik.http.routers.lbi-noauth.rule=Host(`localhost`) && PathPrefix(`/list/`)"
    # Authenticate the rest of the endpoints with basic auth and pass X-Forwarded-User
    - "traefik.http.routers.lbi.rule=Host(`localhost`)"
    - "traefik.http.routers.lbi.middlewares=lbi-auth"
    - "traefik.http.middlewares.lbi-auth.basicauth.headerField=X-Forwarded-User"
    - "traefik.http.middlewares.lbi-auth.basicauth.users=test:$$apr1$$H6uskkkW$$IgXLP6ewTrSuBkTrqE8wj/,test2:$$apr1$$d9hr9HBB$$4HxwgUir3HP4EsggP/QNo0"
```

### OAuth2 authentication with OAuth2 Proxy

This proxy has been tested with letsblockit and can be used with
[a long list of providers](https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/oauth_provider).

- Run the server with
  - `LETSBLOCKIT_AUTH_METHOD=proxy`
  - `LETSBLOCKIT_AUTH_PROXY_HEADER_NAME` set to either `X-Forwarded-Email` or `X-Forwarded-User` depending on your provider
- Set the following options in oauth2-proxy's configuration:

```ini
pass_user_headers = true
skip_auth_routes = [ "GET=^/list/" ]
```

**Warning:** when using an external identity provider, anyone with an account on that platform can login by default.
Check out the per-provider configuration options for details on restricting access to specific users / groups.

### Using a self-hosted Kratos or Ory Cloud

Running with a self-hosted Kratos or even Ory Cloud should work (you'll need to set `LETSBLOCKIT_AUTH_METHOD` to `kratos`
and `LETSBLOCKIT_AUTH_KRATOS_URL`). Don't hesitate to [open an issue](https://github.com/letsblockit/letsblockit/issues/new)
for assistance configuring Kratos itself.
