package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/citihub/probr-sdk/plugin"

	"github.com/citihub/probr-core/internal/core"
)

var packName, varsFile string

func main() {

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
