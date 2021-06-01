package run

import (
	"log"
	"os"
	"os/exec"

	"github.com/probr/probr-sdk/plugin"

	"github.com/probr/probr/internal/core"
)

// CLIContext executes all plugins with handling for the command line
func CLIContext() {
	// Setup for handling SIGTERM (Ctrl+C)
	core.SetupCloseHandler()

	cmdSet, err := core.GetCommands()
	if err != nil {
		log.Printf("Error loading plugins from config: %s", err)
		os.Exit(2)
	}

	// Run all plugins
	if err := AllPlugins(cmdSet); err != nil {
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

// AllPlugins executes specified plugins in a loop
func AllPlugins(cmdSet []*exec.Cmd) (err error) {
	spErrors := make([]core.ServicePackError, 0) // This will store any plugin errors received during execution

	for _, cmd := range cmdSet {
		spErrors, err = Plugin(cmd, spErrors)
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

// Plugin executes single plugin based on the provided command
func Plugin(cmd *exec.Cmd, spErrors []core.ServicePackError) ([]core.ServicePackError, error) {
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
