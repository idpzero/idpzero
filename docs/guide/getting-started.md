---
outline: deep
---

# Getting Started

## Install

Select the installation approach that works for you, ensuring it is added to your path.

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

Once the binary has been installed you can initalize a new ```.idpzero``` directory and related configuration by running the init command. Generally this is only done once per application, and the configuration directory added to source control.

::: tip NOTE
```.idpzero``` has been designed to work nicely with source control. It is recommended that you run the following initialization command within the root of the repository.
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

```sh
Initializing configuration.
(✔) Default configuration initialized.
(✔) Sample configuration applied.

Update your .gitignore to include .idpzero/cache as this directory should not be added to source control.

To start the server run idpzero serve

```

::: warning
The `.idpzero\cache` directory is not intended to be shared and should be added to your `.gitignore`
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

```sh
Running pre-statup checks...
(✔) Verifying cache [OK]
(✔) Checking cache for JWT signing key [OK]
(✔) Checking configuration for clients [OK]
(✔) Checking configuration for users [OK]

Identity Provider started at http://localhost:4379
```

You can shut down `idpzero` by simply ending the running command.

Thats it, your ready to go!