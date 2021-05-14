package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/tabwriter"

	"github.com/probr/probr-sdk/plugin"

	"github.com/probr/probr/internal/core"
	"github.com/probr/probr/internal/flags"
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
		run()
	}
}

func run() {
	// Setup for handling SIGTERM (Ctrl+C)
	core.SetupCloseHandler()

	cmdSet, err := core.GetCommands()
	if err != nil {
		log.Printf("Error loading plugins from config: %s", err)
		os.Exit(2)
	}

	// Run all plugins
	if err := runAllPlugins(cmdSet); err != nil {
		switch e := err.(type) {
		case *core.ServicePackErrors:
			log.Printf("Test Failures: %d out of %d test service packs failed", len(e.SPErrs), len(cmdSet))
			log.Printf("Failed service packs: %v", e.SPErrs)
			os.Exit(1) // At least one service pack failed
		default:
			log.Printf("Internal plugin error: %v", err)
			os.Exit(2) // Internal error
		}
	}
	log.Printf("Success")
	os.Exit(0)
}

func runAllPlugins(cmdSet []*exec.Cmd) (err error) {
	spErrors := make([]core.ServicePackError, 0) // This will store any plugin errors received during execution

	for _, cmd := range cmdSet {
		spErrors, err = runPlugin(cmd, spErrors)
		if err != nil {
			return
		}
	}

	if len(spErrors) > 0 {
		// Return all service pack errors to main
		err = &core.ServicePackErrors{
			SPErrs: spErrors,
		}
	}
	return
}

func runPlugin(cmd *exec.Cmd, spErrors []core.ServicePackError) ([]core.ServicePackError, error) {
	// Launch the plugin process
	client := core.NewClient(cmd)
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return spErrors, err
	}

	// Request the plugin
	rawSP, err := rpcClient.Dispense(plugin.ServicePackPluginName)
	if err != nil {
		return spErrors, err
	}

	// Execute service pack, expecting a silent response
	servicePack := rawSP.(plugin.ServicePack)
	response := servicePack.RunProbes()
	if response != nil {
		spErr := core.ServicePackError{
			ServicePack: cmd.String(), // TODO: retrieve service pack name from interface function
			Err:         response,
		}
		spErrors = append(spErrors, spErr)
	} else {
		log.Printf("[INFO] Probes all completed with successful results")
	}
	return spErrors, nil
}

// listServicePacks lists all service packs declared in config and checks if they are installed
func listServicePacks() {

	servicePackNames, err := core.GetPackNames()
	if err != nil {
		log.Fatalf("An error occurred while retriveing service packs from config: %v", err)
	}

	servicePacks := make(map[string]string)
	for _, pack := range servicePackNames {
		binaryPath, binErr := core.GetPackBinary(pack)
		binaryName := filepath.Base(binaryPath)
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
		fmt.Fprintf(writer, "| %s\t | %s\n", k, v)
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
