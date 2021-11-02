<img src="site/static/images/probr_wide.png">

## Dynamic Application Security Testing (DAST) for Cloud

Probr analyzes the complex behaviours and interactions in your cloud resources to enable engineers, developers and operations teams identify and fix security related flaws at different points in the lifecycle.

Probr has been designed to test aspects of security and compliance that are otherwise challenging to assert using static code inspection or configuration inspection alone. It can also provide a deeper level of confidence in the compliance of your cloud solutions, for those high stakes situations where trusting what your cloud provider is telling you isn't quite enough (software has bugs, after all).

### Control Specifications

Probr uses a structured natural language (Gherkin) to describe the behaviours of an adequately controlled set of cloud resources. These form the basis of control requirements without getting into the nitty gritty of how those controls should be implemented.  This leaves engineering teams the freedom to determine the best course of action to implement the controls that result in those behaviours.

The implementation may change frequently, given the rapid feature velocity in the cloud and tooling ecosystem, without needing to update Probr. This differentiates Probr from policy-based tools, which are designed to look for implementation specifics, so need to iterate in-line with changes to the underlying implementation approach.

### How it works

Probr deploys a series of probes to test the behaviours of the cloud resources in your code, returning a machine-readable set of structured results that can be integrated into the broader DevSecOps process for decision making.  These probes could be as simple as deploying a Kubernetes Pod and running a command inside of it, to complex control and data plane interactions.  If your control can be described as a behaviour then Probr can probe it.

## Architecture

The architecture consists of Probr Core (this repo) and independent service packs containing probes for specific services.  We have built a number of service packs, but you can also build your own using the [Probr SDK](https://github.com/probr/probr-sdk) and following the [boiler plate code](https://github.com/probr/probr-pack-wireframe).

## Available Service Packs

- [Kubernetes core](https://github.com/probr/probr-pack-kubernetes) - cross distribution Kubernetes probes
- [Azure Kubernetes Service (AKS)](https://github.com/probr/probr-pack-aks) - compliments the Kubernetes core pack with AKS specific probes
- [Azure Storage Accounts](https://github.com/probr/probr-pack-storage)

## Quickstart Guide

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

## Configuration

Probr is designed to harness configuration variables to properly target your environment and tailor service pack execution to your unique invironment.

Configuration variables can be populated in multiple ways, with the highest priority value taking precedence.

1. Default values; found in `internal/config/defaults.go` (lowest priority)
1. OS environment variables; set locally prior to probr execution (mid priority)
1. Vars file; See `example-config.yml` in this repository for an example (highest non-CLI priority)

**For more information:** Probr SDK and each service pack use a function named `setEnvAndDefaults` which is used to [wrap the setters for these env vars](https://github.com/probr/probr-sdk/blob/main/config/config.go). By looking directly at this code you can see the names of any config file variable (ex. `ctx.VarName`), the env vars that will be read (ie. `PROBR_VAR_NAME`), and the default value that will be used if neither of the others are provided.

_Note: Different service packs have different requirements, Please see individual service pack documentation for information on the required and default configuations for those packs._

## Special Thanks

We are extremely grateful to the [previous owners](https://github.com/probr-uzh/probr) of this github organization for donating this namespace to our project!
