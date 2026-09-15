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
container image to the GitHub Container Registry. The image is built on Google's
`distroless` base — no shell or package manager — and runs as a non-root user.

The container reads its configuration from `/config`, so mount the directory
that contains your `server.yaml` (and the generated `cache/` state) there.

Initialize a configuration directory (once), then serve it:

```sh
# Create ./idpzero-config/server.yaml with sample clients + users
docker run --rm \
  -v "$(pwd)/idpzero-config:/config" \
  ghcr.io/idpzero/idpzero:latest init --config /config --with-sample-config

# Start the IDP on http://localhost:4379
docker run --rm \
  -p 4379:4379 \
  -v "$(pwd)/idpzero-config:/config" \
  ghcr.io/idpzero/idpzero:latest
```

The default command is `serve --config /config`, so no extra arguments are
needed to start the server.

> **Note:** the image runs as uid/gid `65532`. When bind-mounting a host
> directory the container must be able to write `cache/state.sqlite` into it.
> If you hit permission errors, run the container as your own user by adding
> `--user "$(id -u):$(id -g)"` to the commands above (a named volume works
> without this).

To build the image from source instead of pulling a release:

```sh
make docker            # builds idpzero:local
# or: docker build -t idpzero:local .
```

