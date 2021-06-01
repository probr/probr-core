package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/probr/probr/internal/core"
	"github.com/probr/probr/internal/flags"
	"github.com/probr/probr/run"
)

var (
	// See Makefile for more on how this package is built

	// Version is the main version number that is being run at the moment
	Version = "0.1.0"

	// VersionPostfix is a marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc", etc.
	VersionPostfix = "dev"

	// GitCommitHash references the commit id at build time
	GitCommitHash = ""

	// BuiltAt is the build date
	BuiltAt = ""
)

func main() {

	var subCommand string
	if len(os.Args) > 1 {
		subCommand = os.Args[1]
	}
	switch subCommand {
	// Ref: https://gobyexample.com/command-line-subcommands
	case "list":
		flags.List.Parse(os.Args[2:])
		listServicePacks()

	case "version":
		flags.Version.Parse(os.Args[2:])
		printVersion()

	default:
		flags.Run.Parse(os.Args[1:])
		run.CLIContext()
	}
}

func printVersion() {
	if VersionPostfix != "" {
		Version = fmt.Sprintf("%s-%s", Version, VersionPostfix)
	}

	fmt.Fprintf(os.Stdout, "Probr Version: %s", Version)
	if core.Verbose != nil && *core.Verbose {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Commit       : %s", GitCommitHash)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Built at     : %s", BuiltAt)
	}
}

// listServicePacks reads all service packs declared in config and checks whether they are installed
func listServicePacks() {
	servicePackNames, err := core.GetPackNames()
	if err != nil {
		log.Fatalf("An error occurred while retrieving service packs from config: %v", err)
	}

	servicePacks := make(map[string]string)
	for _, pack := range servicePackNames {
		_, binErr := core.GetPackBinary(pack)
		if binErr != nil {
			servicePacks[pack] = fmt.Sprintf("ERROR: %v", binErr)
		} else {
			servicePacks[pack] = "OK"
		}
	}

	// Print output
	writer := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(writer, "| Service Pack\t | Installed ")
	for k, v := range servicePacks {
		fmt.Fprintf(writer, "| %s\t | %s\n", k, v)
	}
	writer.Flush()
}
