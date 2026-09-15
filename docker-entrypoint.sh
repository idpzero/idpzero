#!/bin/sh
set -e

# idpzero stores its state database and generated JWT signing key under the
# configuration directory, so that directory must be writable at runtime.
#
# This entrypoint follows the same pattern as the official Postgres image: when
# the container starts as root (the default), fix ownership of the mounted
# configuration directory and then drop to the unprivileged "idpzero" user
# before running the server. This lets a plain `docker run -v ./dir:/config`
# work without the caller having to match uids.
#
# If the container was started as a non-root user (e.g. `docker run --user ...`
# or a Kubernetes securityContext), we skip the chown and exec directly, trusting
# that the platform has made the mount writable.

CONFIG_DIR="${IDPZERO_CONFIG_DIR:-/config}"

if [ "$(id -u)" = "0" ]; then
    if [ -d "$CONFIG_DIR" ]; then
        chown -R idpzero:idpzero "$CONFIG_DIR" 2>/dev/null || \
            echo "idpzero: warning: could not change ownership of $CONFIG_DIR (read-only mount?)" >&2
    fi
    exec su-exec idpzero:idpzero /usr/bin/idpzero "$@"
fi

exec /usr/bin/idpzero "$@"
