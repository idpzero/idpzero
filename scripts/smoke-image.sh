#!/bin/sh
# Smoke test a built idpzero image: serve a minimal configuration and check
# that discovery and the logged-out page answer.
#
#   scripts/smoke-image.sh [image]   (default idpzero:local, as built by make docker)
set -eu

image="${1:-idpzero:local}"
port="${SMOKE_PORT:-14379}"
name="idpzero-smoke-$$"
dir="$(mktemp -d)"

# Only config/ is mounted: the entrypoint chowns it to the image's idpzero
# user, after which (on Linux) this user can no longer write there or delete
# it, so outputs live beside it and the container removes it.
mkdir "$dir/config"

cleanup() {
    docker rm -f "$name" >/dev/null 2>&1 || true
    docker run --rm --user 0 --entrypoint rm -v "$dir:/smoke" "$image" -rf /smoke/config >/dev/null 2>&1 || true
    rm -rf "$dir" 2>/dev/null || true
}
trap cleanup EXIT

cat > "$dir/config/server.yaml" <<YAML
server:
  port: 4379
  keyphrase: smoke-test
clients:
  - name: Smoke
    client_id: smoke
    grant_types: [authorization_code]
    response_types: [code]
    redirect_uris: [http://localhost/callback]
users:
  - subject: alice
    login_display: Alice
YAML

docker run -d --name "$name" -p "$port:4379" -v "$dir/config:/config" "$image" >/dev/null

base="http://localhost:$port"
i=0
until curl -fsS "$base/.well-known/openid-configuration" -o "$dir/discovery.json" 2>/dev/null; do
    i=$((i + 1))
    if [ "$i" -ge 60 ]; then
        echo "smoke: discovery did not answer within 30s" >&2
        curl -sS "$base/.well-known/openid-configuration" -o /dev/null >&2 || true
        docker logs "$name" >&2
        exit 1
    fi
    sleep 0.5
done

if ! grep -q "\"issuer\":\"$base\"" "$dir/discovery.json"; then
    echo "smoke: unexpected discovery document:" >&2
    cat "$dir/discovery.json" >&2
    exit 1
fi
curl -fsS -o /dev/null "$base/logged-out"

echo "smoke: ok ($image)"
