<img src="assets/probr.png" width="200">

## Dynamic Application Security Testing (DAST) for Cloud
Probr analyzes the complex behaviours and interactions in your cloud resources to enable engineers, developers and operations teams identify and fix security related flaws at different points in the lifecycle.

Probr has been designed to test aspects of security and compliance that are otherwise challenging to assert using static code inspection or configuration inspection alone. It can also provide a deeper level of confidence in the compliance of your cloud solutions, for those high stakes situations where trusting what your cloud provider is telling you isn't quite enough (software has bugs, after all).

### Control Specifications
Probr uses a structured natural language (Gherkin) to describe the behaviours of an adequately controlled set of cloud resources. These form the basis of control requirements without getting into the nitty gritty of how those controls should be implemented.  This leaves engineering teams the freedom to determine the best course of action to implement the controls that result in those behaviours.

The implementation may change frequently, given the rapid feature velocity in the cloud and tooling ecosystem, without needing to update Probr. This differentiates Probr from policy-based tools, which are designed to look for implementation specifics, so need to iterate in-line with changes to the underlying implementation approach.

### How it works
Probr deploys a series of probes to test the behaviours of the cloud resources in your code, returning a machine-readable set of structured results that can be integrated into the broader DevSecOps process for decision making.  These probes could be as simple as deploying a Kubernetes Pod and running a command inside of it, to complex control and data plane interactions.  If your control can be described as a behaviour then Probr can probe it.

## Architecture

The architecture consists of Probr Core (this repo) and independent service packs containing probes for specific services.  We have built a number of service packs, but you can also build your own using the [Probr SDK](https://github.com/citihub/probr-sdk).  We have a developer guide and boiler plate code here (to be done).

## Available Service Packs

- [Kubernetes core](https://github.com/citihub/probr-pack-kubernetes) - cross distribution Kubernetes probes
- [Azure Kubernetes Service (AKS)](https://github.com/citihub/probr-pack-aks) - compliments the Kubernetes core pack with AKS specific probes
- [Azure Storage Accounts](https://github.com/citihub/probr-pack-storage)

## Quickstart Guide

### Get the Probr executable

- **Option 1** - Download the latest Probr package by clicking the corresponding asset on our [release page](https://github.com/citihub/probr-core/releases).
- **Option 2** - You may build the edge version of Probr by using `make binary` from the source code. This may also be necessary if an executable compatible with your system is not available in on the release page.
- **Option 3** - TODO: Example Dockerfile which will build a Docker image with both Probr and [Cucumber HTML Reporter](https://www.npmjs.com/package/cucumber-html-reporter) for visualisation

*Note: The usage docs refer to the executable as `probr` or `probr.exe` interchangeably. Use the former for unix/linux systems, and the latter package if you are working in Windows.*

### Get a service pack

See individual service packs for instructions on how to obtain the binary.

By default Probr will look in the `${HOME}/probr/binaries` path for the service packs. If you want to put them in a different location then you can use the `-binaries-path <directory>` flag when running Probr.

### Configure Probr

Configuration variables can be populated in one of four ways, with the value being taken from the highest priority entry.

1. Default values; found in `internal/config/defaults.go` (lowest priority)
1. OS environment variables; set locally prior to probr execution (mid priority)
1. Vars file; yaml (highest non-CLI priority)
1. CLI flags; see `./probr --help` for available flags (highest priority)

See `config.yml` in this repository for an example of configuring Probr.  If you just want to try it out then the defaults will usually be sufficient.

_Note: Different service packs have different requirements, Please see individual service pack documentation for information on the required and default configuations for those packs._

### Run the CLI

1. Run the probr executable via `./probr [OPTIONS]`.  By default it will look for `config.yml` in the same location that you run probr from.
    - If your binaries aren't in `${HOME}/probr/binaries` then use `-binaries-path=<path>`.
    - Other options can be seen via `./probr --help`

### View the results

The default location for Probr output is `${HOME}/probr/output/<date>/<time>/<service_pack>`. There are various output files, as follows...

#### Summary results

`summary.json` displays an overall summary of the Probr results.

#### Cucumber results

In the `cucumber` sub-folder the Probr results are displayed in a standard "Cucumber" JSON format, which can be fed into your favourite Cucumber parser or visualisation tool.

#### Audit trail

In the `audit` sub-folder, there is an audit trail of every step the service pack executed in deploying the probe.  For example, the Kubernetes service pack audit trail captures the exact pod specifications that were deployed for each probe and the response received from Kubernetes.

## More configuration

### Environment Variables

If you would like to handle logic differently per environment, env vars may be useful. An example of how to set an env var is as follows:

`export KUBE_CONFIG=./path/to/config`

### Vars File

An example Vars file is available at [./examples/config.yml](./examples/config.yml).
You may have as many vars files as you wish in your codebase, which will enable you to maintain configurations for multiple environments in a single codebase.

The location of the vars file is passed as a CLI option e.g.

```
./probr --varsFile=./config-dev.yml
```

### Probr Configuration Variables

These are general configuration variables.

| Variable | Description | CLI Option | Vars File | Env Var | Default |
|---|---|---|---|---|---|
|VarsFile|Config YAML File Path|yes|N/A|N/A|N/A|
|Silent|Disable visual runtime indicator|yes|no|N/A|false|
|NoSummary|Flag to switch off summary output|yes|no|N/A|false|
|WriteDirectory|Path to all output, including audit, cucumber results and other temp files|yes|yes|PROBR_WRITE_DIRECTORY|probr_output|
|Tags|Feature tag inclusions and exclusions|yes|yes|PROBR_TAGS| |
|LogLevel|Set log verbosity level|yes|yes|PROBR_LOG_LEVEL|ERROR|
|OutputType|"IO" will write to file, as is needed for CLI usage. "INMEM" should be used in non-CLI cases, where values should be returned in-memory instead|no|yes|PROBR_OUTPUT_TYPE|IO|
|AuditEnabled|Flag to switch on audit log|no|yes|PROBR_AUDIT_ENABLED|true|
|OverwriteHistoricalAudits|Flag to allow audit overwriting|no|yes|OVERWRITE_AUDITS|true|

## Development & Contributing

Please see the [contributing docs](https://github.com/citihub/probr/blob/master/CONTRIBUTING.md) for information on how to develop and contribute to this repository as either a maintainer or open source contributor (the same rules apply for both).
