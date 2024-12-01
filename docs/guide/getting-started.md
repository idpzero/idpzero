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
