# syntax=docker/dockerfile:1
#
# Hardened, multi-stage build of the idpzero binary.
#
# The runtime image is Google's "distroless static" base: it contains only the
# binary, CA certificates and tzdata. There is no shell, package manager or
# libc, and it runs as an unprivileged user (uid/gid 65532) by default, which
# keeps the attack surface minimal.
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

# Pre-create the config/state directory so named volumes inherit non-root
# ownership (distroless has no shell to mkdir/chown at runtime).
RUN mkdir -p /out/config

# ---- Runtime stage ---------------------------------------------------------
FROM gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.title="idpzero" \
      org.opencontainers.image.description="Single binary Identity Provider (IDP) for local dev/test" \
      org.opencontainers.image.source="https://github.com/idpzero/idpzero" \
      org.opencontainers.image.url="https://idpzero.dev" \
      org.opencontainers.image.licenses="MIT"

COPY --from=build /out/idpzero /usr/bin/idpzero
COPY --from=build --chown=65532:65532 /out/config /config

# Configuration + state directory. Bind-mount your host `.idpzero` directory
# here (or use a named volume). idpzero reads server.yaml and writes
# cache/state.sqlite under this path.
VOLUME ["/config"]

# Default IDP listen port (server.port in server.yaml).
EXPOSE 4379

USER 65532:65532

ENTRYPOINT ["/usr/bin/idpzero"]
CMD ["serve", "--config", "/config"]
