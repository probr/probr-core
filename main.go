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
	// Ref: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
	// Below are some examples for setting version during build time. This could be used in a make file and/or in CI/CD pipeline (preferred).
	// Local dev:
	//   > go build -o probr -ldflags="-X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`'"
	// Release candidate:
	//   > go build -o probr -ldflags="-X 'main.VersionPostfix=rc' -X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`'"
	// Production release:
	//   > go build -o probr -ldflags="-X 'main.VersionPostfix=' -X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`'"
	// Setting all version details inline:
	//   > go build -o probr -ldflags="-X 'main.Version=0.14.0' -X 'main.VersionPostfix=rc' -X 'main.GitCommitHash=`git rev-parse --short HEAD`' -X 'main.BuiltAt=`date +%FT%T%z`'"

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

var packName, varsFile string

func main() {

	// > probr list [-binaries-path]
	coreCmd := flag.NewFlagSet("probr", flag.ExitOnError)
	core.BinariesPath = coreCmd.String("binaries-path", "", "Location for service pack binaries. If not provided, default value is: [UserHomeDir]/probr/binaries")

	// > probr version [-v]
	versionCmd := flag.NewFlagSet("probr version", flag.ExitOnError)
	verboseVersionFlag := versionCmd.Bool("v", false, "Display extended version information")

	subCommand := ""
	if len(os.Args) > 1 {
		subCommand = os.Args[1]
	}
	switch subCommand {
	case "list":
		coreCmd.Parse(os.Args[2:])
		listServicePacks(os.Stdout)

	case "version":
		versionCmd.Parse(os.Args[2:])
		printVersion(os.Stdout, *verboseVersionFlag)

	default:
		coreCmd.Parse(os.Args[1:])
		runServicePacks()
	}
	// Ref for handling cli subcommands: https://gobyexample.com/command-line-subcommands
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

func printVersion(w io.Writer, verbose bool) {

	fmt.Fprintf(w, "Probr Version: %s", getVersion())
	if verbose {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Commit       : %s", GitCommitHash)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Built at     : %s", BuiltAt)
	}
}

func getVersion() string {
	if VersionPostfix != "" {
		return fmt.Sprintf("%s-%s", Version, VersionPostfix)
	}
	return Version
}
