package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/tabwriter"

	"github.com/citihub/probr-sdk/plugin"

	"github.com/citihub/probr-core/internal/core"
	"github.com/citihub/probr-core/internal/flags"
)

var (
	// See Makefile for more on how this package is built

	// Version is the main version number that is being run at the moment
	Version = "0.0.15"

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
		flags.Core.Parse(os.Args[1:])
		runServicePacks()
	}
}

func runServicePacks() {
	// Setup for handling SIGTERM (Ctrl+C)
	core.SetupCloseHandler()

	cmdSet, err := core.GetCommands()
	if err != nil {
		log.Printf("Error loading plugins from config: %s", err)
		os.Exit(2)
	}

	totalCount := len(cmdSet)
	if totalCount == 0 {
		log.Print("No service pack found in config")
		return
	}

	// Run all plugins
	//  exit 2 on internal error
	//  exit 1 on service pack error(s)
	//  exit 0 on success
	if err := runAllPlugins(cmdSet); err != nil {
		switch e := err.(type) {
		case *core.ServicePackErrors:
			log.Printf("Test Failures: %d out of %d test service packs failed", len(e.SPErrs), totalCount)
			log.Printf("Failed service packs: %v", e.SPErrs)
			os.Exit(1) // At least one service pack failed
		default:
			log.Printf("Internal plugin error: %v", err)
			os.Exit(2) // Internal error
		}
	}
}

func runAllPlugins(cmdSet []*exec.Cmd) error {
	var err error
	spErrors := make([]core.ServicePackError, 0) // Intialize collection to store any service pack error received during plugin execution

	for _, cmd := range cmdSet {
		// Launch the plugin process
		client := core.NewClient(cmd)
		defer client.Kill()

		// Connect via RPC
		rpcClient, err := client.Client()
		if err != nil {
			return err
		}

		// Request the plugin
		rawSP, err := rpcClient.Dispense(plugin.ServicePackPluginName)
		if err != nil {
			return err
		}
		// We should have a ServicePack now! This feels like a normal interface
		// implementation but is in fact over an RPC connection.
		servicePack := rawSP.(plugin.ServicePack)
		result := servicePack.RunProbes()
		if result != nil {
			spErr := &core.ServicePackError{
				ServicePack: cmd.String(),
				Err:         result,
			}
			spErrors = append(spErrors, *spErr)
		} else {
			log.Printf("[INFO] Probes all completed with successful results")
		}
		// Confirmed this can handled long-running plugins
		// It worked with simulated 30-second delay
		// It worked with simulated 10-min delay
		// It worked with simulated 30-min delay
	}

	if len(spErrors) > 0 {
		// Return all service pack errors to main
		err = &core.ServicePackErrors{
			SPErrs: spErrors,
		}
	}

	return err
}

//listServicePacks lists all service packs declared in config and checks if they are installed
func listServicePacks() {

	servicePackNames, err := core.GetPackNames()
	if err != nil {
		log.Fatalf("An error occurred while retriveing service packs from config: %v", err)
	}

	servicePacks := make(map[string]string)
	for _, pack := range servicePackNames {
		binaryName, binErr := core.GetPackBinary(pack)
		if binErr != nil {
			servicePacks[binaryName] = fmt.Sprintf("ERROR: %v", binErr)
		} else {
			servicePacks[binaryName] = "OK"
		}
	}

	// Print output
	writer := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(writer, "| Service Pack\t | Installed ")
	for k, v := range servicePacks {
		fmt.Fprintln(writer, fmt.Sprintf("| %s\t | %s", k, v))
	}
	writer.Flush()
}

func printVersion() {
	fmt.Fprintf(os.Stdout, "Probr Version: %s", getVersion())
	if core.Verbose != nil && *core.Verbose {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Commit       : %s", GitCommitHash)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "Built at     : %s", BuiltAt)
	}
}

func getVersion() string {
	if VersionPostfix != "" {
		return fmt.Sprintf("%s-%s", Version, VersionPostfix)
	}
	return Version
}
