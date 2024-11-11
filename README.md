# idpzero
Local IdP for development and testing purposes, single binary with zero addition dependencies


TODO:
 - https://templ.guide/commands-and-tools/live-reload-with-other-tools/

 Run Sample Client:
 ```
 CLIENT_ID=web CLIENT_SECRET=secret ISSUER=http://localhost:4379 SCOPES="openid profile" PORT=9999 go run github.com/zitadel/oidc/v3/example/client/app
 ```

 Navigate to http://localhost:9999/login
 
