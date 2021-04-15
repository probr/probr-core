package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/citihub/probr-sdk/plugin"

	"github.com/citihub/probr-core/internal/core"
)

var (
	// Version is the main version number that is being run at the moment
	Version = "0.0.0"

	// Prerelease is a marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc1", etc.
	Prerelease = "dev"

	// GitCommitHash shall be used to store commit id when building release
	GitCommitHash = ""
)

var packName, varsFile string

func main() {

	// > probr list
	flag.NewFlagSet("list", flag.ExitOnError)

	// > probr version
	flag.NewFlagSet("version", flag.ExitOnError)

	subCommand := ""
	if len(os.Args) > 1 {
		subCommand = os.Args[1]
	}
	switch subCommand {
	case "list":
		listServicePacks(os.Stdout)

	case "version":
		printVersion(os.Stdout)

	default:
		runServicePacks()
	}
	// Ref for handling cli subcommands: https://gobyexample.com/command-line-subcommands
}

func runServicePacks() {
	// Setup for handling SIGTERM (Ctrl+C)
	core.SetupCloseHandler()

	core.ParseFlags()
	// if err != nil {
	// 	log.Printf("Error parsing flags from command line: %s", err)
	// 	os.Exit(2)
	// }

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
		// TODO: Confirm how to handle long-running plugins, since this will block until finished. Potential time out.
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
func listServicePacks(w io.Writer) {

	declaredServicePacks, err := core.GetPackNameFromConfig()
	if err != nil {
		log.Fatalf("An error occurred while retriveing service packs from config: %v", err)
	}

	servicePacks := make(map[string]string)

	for _, pack := range declaredServicePacks {
		_, binErr := core.GetPackBinary(pack)
		if binErr != nil {
			servicePacks[pack] = fmt.Sprintf("ERROR: %v", binErr)
		} else {
			servicePacks[pack] = "OK"
		}
	}

	// Print output
	fmt.Fprintln(w, "Listing all declared service packs... ")
	fmt.Fprintln(w, "| Service Pack\t\t\t\t | Installed ")
	for k, v := range servicePacks {
		fmt.Fprintln(w, fmt.Sprintf("| %s\t\t\t\t | %s", k, v))
	}
}

func printVersion(w io.Writer) {

	fmt.Fprintf(w, "Probr Version: %s", getVersion())
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Commit: %s", GitCommitHash)
	//Ref: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
	// To set a version during build time: go build -o probr -ldflags="-X 'main.Version=0.12.0' -X 'main.Prerelease=rc' -X 'main.GitCommitHash=123456'"
	// This could be used in a make file and/or in CI/CD pipeline (preferred)
}

func getVersion() string {
	if Prerelease != "" {
		return fmt.Sprintf("%s-%s", Version, Prerelease)
	}
	return Version
}
