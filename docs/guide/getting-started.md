---
outline: deep
---

# Getting Started

## Install

Select the installation approach that works for you, or download from Github [releases](https://github.com/idpzero/idpzero/releases). Make sure the binary is available in your local path.

::: code-group

```sh [go]
$ go install http://github.com/idpzero/idpzero
```

```sh [brew]
$ brew tap idpzero/idpzero
$ brew install idpzero
```
:::


## Initialize

Once the binary has been installed you can initalize a new configuration directory by running the `init` command. Generally this is only done once per application, and the configuration directory added to source control.

::: tip NOTE
The generate `.idpzero` directory has been designed to work nicely with source control. It is recommended that you run the following initialization command within the root of the repository.
:::

The following includes the `--with-sample-config` flag which will generate sample client and user configuration to help get started. 

::: code-group
```sh [linux/mac]
# Initialize the .idpzero folder to current directory.
# Add --help flag for more options.
idpzero init --with-sample-config
```
```sh [windows]
# Initialize the .idpzero folder to current directory.
# Add --help flag for more options.
idpzero.exe init --with-sample-config
```

:::

You should see output similar to the below

![idpzero init](/shell/init.gif)

::: warning
The `.idpzero/cache` directory should be added to your `.gitignore`
:::


If you explore the `.idpzero` directory that has been created, you will see a `server.yaml` file which contains the configuration.


## Start the IDP

Starting the IDP is as simple as running the `serve` command within the CLI. The default port that the server will be served on is `4379` unless overriden within the ```.idpzero``` configuration directory.

::: tip NOTE
By default `idpzero` will auto discover the closest ```.idpzero``` configuration folder by walking up the directory tree from the current location. See `--help` for additional options.
:::

::: code-group
```sh [linux/mac]
# Start idpzero using configuration discovery
# Add --help flag for more options.
idpzero serve
```
```sh [windows]
# Start idpzero using configuration discovery
# Add --help flag for more options.
idpzero.exe serve
```

:::

As part of the startup, various checks will be executed, and you will see output similar to below.

![idpzero init](/shell/serve.gif)

You can shut down server by simply ending the running command.

Thats it, your ready to go!

## View Configuration

Once running, open the dashboard (default is http://localhost:4379/) where configuration for the server including metadata endpoints to add to your application, as well as configured client information can be viewed.

If client secrets are applicable for the application, these will be available here to copy.

![Dashboard](/screenshots/dashboard.png)

## Client Authentication (PKCE vs Client Secret)

`idpzero` supports both public and confidential OAuth2 / OpenID Connect clients. The
`auth_method` field on each client controls how the client authenticates when it
exchanges an authorization code for tokens.

### Public clients (PKCE — no secret required)

For single page apps (SPAs), native/mobile apps, CLIs and anything else that cannot
keep a secret, use a **public client** authenticating with
[PKCE](https://datatracker.ietf.org/doc/html/rfc7636). Set `auth_method: none`:

```yaml
clients:
  - name: My SPA
    client_id: spa
    application_type: native
    auth_method: none          # public client, authenticates via PKCE
    grant_types:
      - authorization_code
      - refresh_token
    redirect_uris:
      - http://localhost:3000/callback
    response_types:
      - code
```

::: tip NO SECRETS TO STORE
Public clients have **no client secret**. Because there is nothing sensitive to
store, this configuration can be committed to source control as-is. The
authorization code exchange is protected by the PKCE `code_verifier` /
`code_challenge` (`S256`) instead of a shared secret.
:::

If `auth_method` is omitted it defaults to `none`, so the safest, secretless
public-client behaviour is the default.

### Confidential clients (client secret)

For server-side web apps that can safely store a secret, use
`auth_method: client_secret_basic` (or `client_secret_post`). `idpzero` derives a
stable client secret from the server `keyphrase` and the `client_id`, and displays
it on the dashboard so you can copy it into your application configuration:

```yaml
clients:
  - name: My Web App
    client_id: web
    application_type: web
    auth_method: client_secret_basic
    grant_types:
      - authorization_code
      - refresh_token
    redirect_uris:
      - http://localhost:3000/callback
    response_types:
      - code
```

::: warning
The derived secret depends on the `keyphrase` in `server.yaml`. Keep the `keyphrase`
out of source control if you rely on confidential clients; public (PKCE) clients
avoid this concern entirely.
:::

The `--with-sample-config` initializer includes one of each client type to get you
started.

## Provisioning users to another system (SCIM)

`idpzero` can push its configured users to an external SCIM 2.0 service from the
**Provisioning** page in the dashboard. See [SCIM Provisioning](/guide/scim-provisioning)
for configuration and usage.
