# syntax=docker/dockerfile:1
#
# Hardened, multi-stage build of the idpzero binary.
#
# The runtime image is a minimal Alpine base carrying only the static binary,
# CA certificates and su-exec. An entrypoint fixes ownership of the mounted
# config directory and then drops to the unprivileged "idpzero" user before
# running the server (the same drop-privileges pattern the official Postgres
# image uses), so bind-mounted volumes are writable without the caller needing
# to match uids.
#
# Build locally with:
#   docker build -t idpzero:local .
#
# Released images are built by GoReleaser from the pre-compiled binaries (see
# Dockerfile.goreleaser); this file is the equivalent build-from-source path.

# ---- Build stage -----------------------------------------------------------
FROM golang:1.25 AS build

WORKDIR /src

# Download modules first so they are cached independently of source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Version metadata injected into the binary (matches GoReleaser's ldflags).
ARG VERSION=dev
ARG COMMIT=none

# CGO is disabled so the result is a fully static binary with no libc
# dependency. modernc.org/sqlite is a pure-Go driver, so the state database
# needs no cgo either.
RUN CGO_ENABLED=0 GOOS=linux go build \
      -trimpath \
      -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
      -o /out/idpzero .

# ---- Runtime stage ---------------------------------------------------------
FROM alpine:3.21

LABEL org.opencontainers.image.title="idpzero" \
      org.opencontainers.image.description="Single binary Identity Provider (IDP) for local dev/test" \
      org.opencontainers.image.source="https://github.com/idpzero/idpzero" \
      org.opencontainers.image.url="https://idpzero.dev" \
      org.opencontainers.image.licenses="MIT"

# ca-certificates: needed for outbound TLS (e.g. SCIM provisioning).
# su-exec: drop privileges from the entrypoint after fixing volume ownership.
# Create the unprivileged runtime user and pre-create the config dir so named
# volumes inherit its ownership.
RUN apk add --no-cache ca-certificates su-exec \
    && addgroup -g 65532 -S idpzero \
    && adduser -u 65532 -S -H -D -G idpzero idpzero \
    && mkdir -p /config \
    && chown idpzero:idpzero /config

COPY --from=build /out/idpzero /usr/bin/idpzero
COPY --chmod=0755 docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

# Configuration + state directory. Bind-mount your host `.idpzero` directory
# here (or use a named volume). idpzero reads server.yaml and writes
# cache/state.sqlite under this path.
VOLUME ["/config"]

# Default IDP listen port (server.port in server.yaml).
EXPOSE 4379

# The entrypoint starts as root, chowns /config, then drops to the idpzero user
# (uid/gid 65532). Pass `--user` to skip that and run as an arbitrary uid.
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["serve", "--config", "/config"]
