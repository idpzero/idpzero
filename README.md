# idpzero
Local IdP for development and testing purposes, single binary with zero addition dependencies


For local development, we suggest leveraging `make dev` which will setup the required watches across the development tooling required. This will launch a brower at `http://localhost:8080` which operates as a proxy to the `http://localhost:4379` to enable hot reloading.


 Run Sample Client:
 ```
 CLIENT_ID=web CLIENT_SECRET=secret ISSUER=http://localhost:4379 SCOPES="openid profile" PORT=9999 go run github.com/zitadel/oidc/v3/example/client/app
 ```

 Navigate to http://localhost:9999/login

## Running with Docker

Released versions are published as a hardened, multi-arch (amd64/arm64)
container image to the GitHub Container Registry
(`ghcr.io/idpzero/idpzero`). The image runs the server as an unprivileged user
and carries only the static binary, CA certificates and `su-exec` on a minimal
Alpine base.

The container reads its configuration from `/config`, so bind-mount the
directory that holds your `server.yaml` (the generated `cache/` state is written
alongside it). Initialize a configuration directory once, then serve it:

```sh
# Create ./.idpzero/server.yaml with sample clients + users
docker run --rm \
  -v "$(pwd)/.idpzero:/config" \
  ghcr.io/idpzero/idpzero:latest init --config /config --with-sample-config

# Start the IDP on http://localhost:4379
docker run --rm \
  -p 4379:4379 \
  -v "$(pwd)/.idpzero:/config" \
  ghcr.io/idpzero/idpzero:latest
```

The default command is `serve --config /config`, so no extra arguments are
needed to start the server.

> **How volume permissions work:** the container starts as root only long enough
> for its entrypoint to take ownership of `/config`, then drops to the
> unprivileged `idpzero` user (uid/gid `65532`) to run the server — the same
> approach the official Postgres image uses. This means a bind-mounted host
> directory just works, with no `--user` juggling. If you prefer to run as a
> fixed user (e.g. in a locked-down/Kubernetes environment), pass
> `--user "$(id -u):$(id -g)"`; the entrypoint detects it isn't root, skips the
> chown, and runs directly — just ensure the mounted directory is writable by
> that user.

To build the image from source instead of pulling a release:

```sh
make docker            # builds idpzero:local
# or: docker build -t idpzero:local .
```

