# idpzero
Local IdP for development and testing purposes, single binary with zero addition dependencies


For local development, we suggest leveraging `make dev` which will setup the required watches across the development tooling required. This will launch a brower at `http://localhost:8080` which operates as a proxy to the `http://localhost:4379` to enable hot reloading.


 Run Sample Client:
 ```
 CLIENT_ID=web CLIENT_SECRET=secret ISSUER=http://localhost:4379 SCOPES="openid profile" PORT=9999 go run github.com/zitadel/oidc/v3/example/client/app
 ```

 Navigate to http://localhost:9999/login
 
