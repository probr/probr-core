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
