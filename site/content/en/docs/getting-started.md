---
title: "Getting Started"
date: 2021-11-02T08:39:09-05:00
draft: false
---

This guide will help you install and make a basic execution of Probr.

### Get the Probr executable

- **Option 1** - Download the latest Probr package by clicking the corresponding asset on our [release page](https://github.com/probr/probr/releases).
- **Option 2** - You may build the edge version of Probr by using `make binary` from the source code. This may also be necessary if an executable compatible with your system is not available in on the release page.
- **Option 3** - We've [containerized Probr](https://github.com/probr/probr-docker) along with all approved service packs. Visit the repo for more information about how you can harness it for your organization.

### Get a service pack

Each service pack can be retrieved as a binary on the releases page of it's repo, or built using the provided Makefile. Installing the service pack is simply a matter of moving it to the `bin/` directory in your install path.

By default Probr will look in `<HOME>/probr/bin` for the service packs, but if you modify the installation directory (and specify it in your configuration) then Probr will look in the corresponding location: `<INSTALL_DIR>/bin`

### Run the CLI

1. Run the probr executable via `./probr [OPTIONS]`.  By default it will look for `config.yml` in the same location that you run probr from.
    - Other options can be seen via `./probr --help`

### View the results

The default location for Probr output is `${HOME}/probr/output`. There are various output files, as follows...

- `summary.json`
displays an overall summary of the Probr results.
- `cucumber/`
Files containing Probe results are displayed in a standardized format, which can be fed into your favourite Cucumber parser or visualisation tool.
- `audit/`
Files for each probe contain an audit trail of every step that was executed.

### Runtime Configuration

Probr is designed to harness configuration variables to properly target your environment and tailor service pack execution to your unique invironment.

Configuration variables can be populated in multiple ways, with the highest priority value taking precedence.

1. Default values; found in `internal/config/defaults.go` (lowest priority)
1. OS environment variables; set locally prior to probr execution (mid priority)
1. Vars file; See `example-config.yml` in this repository for an example (highest priority)

> **For more information:** Probr SDK and each service pack use a function named `setEnvAndDefaults` which is used to [wrap the setters for these env vars](https://github.com/probr/probr-sdk/blob/main/config/config.go). By looking directly at this code you can see the names of any config file variable (ex. `ctx.VarName`), the env vars that will be read (ie. `PROBR_VAR_NAME`), and the default value that will be used if neither of the others are provided.

The following config executes the only Kubernetes service pack, and passes multiple variables to that service pack.

```
Run:
  - kubernetes
ServicePacks:
  Kubernetes:
    KubeConfig: /probr/run/kubeconfig
    AuthorisedContainerImage: citihubprod.azurecr.io/citihub/probr-probe
    UnauthorisedContainerImage: docker.io/citihub/probr-probe
```

To run multiple service packs, simply add another item to the list of packs to run:

```
Run:
  - kubernetes
  - storage
```

If each pack calls for config variables, nest the vars within the name of the service pack.

```
ServicePacks:
  Kubernetes:
    KubeConfig: /probr/run/kubeconfig
  Storage:
    Key: /probr/run/keyfile
```

_Note: Different service packs have different requirements, Please see individual service pack documentation for information on the required and default configuations for those packs._
