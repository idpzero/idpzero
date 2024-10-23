# idpzero
Local IdP for development and testing purposes, single binary with zero addition dependencies


https://www.ory.sh/docs/hydra/jwks

Example CLI structures

```
ory create jwks \
  --alg RS256 \
  hydra.openid.id-token
```

```
  ory patch oauth2-config --project <project-id> --workspace <workspace-id> \
  --add '/webfinger/jwks/broadcast_keys/-="custom_keyset"'
  ```