package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/plugin"
	hcplugin "github.com/hashicorp/go-plugin"

	"github.com/citihub/probr-core/internal/core"
)

var packName, varsFile string

func main() {

	// Setup for handling SIGTERM (Ctrl+C)
	core.SetupCloseHandler()

	// Handle cli args and load plugins from config file
	cmdSet, err := parseFlags()
	if err != nil {
		log.Printf("Error loading plugins from config: %v", err)
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

func userHomeDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return user.HomeDir
}

func packBinary(name string) (binaryName string, err error) {
	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(name), ".exe") {
		name = fmt.Sprintf("%s.exe", name)
	}
	binaryPath := filepath.Join(userHomeDir(), "probr", "binaries") // TODO Load from config.
	plugins, _ := hcplugin.Discover(name, binaryPath)
	if len(plugins) != 1 {
		err = fmt.Errorf("Please ensure requested plugin '%s' has been installed to '%s'", name, binaryPath)
		return
	}
	binaryName = plugins[0]

	return
}

func parseFlags() (cmdSet []*exec.Cmd, err error) {
	var configPath string
	argCount := len(os.Args)
	if argCount < 2 {
		err = errors.New("First argument should path to config file")
		return
	}
	configPath = os.Args[1]

	packNames, err := getPackNameFromConfig(configPath)
	if err != nil {
		return
	}

	for _, pack := range packNames {
		binaryName, binErr := packBinary(pack)
		if binErr != nil {
			err = binErr
			break
		}
		cmd := exec.Command(binaryName)
		cmd.Args = append(cmd.Args, fmt.Sprintf("--varsfile=%s", configPath))

		if argCount > 2 {
			// TODO: passing flags to service pack isn't scalable
			cmd.Args = append(cmd.Args, os.Args[2:]...)
		}
		cmdSet = append(cmdSet, cmd)
	}

	return
}

func getPackNameFromConfig(configPath string) (packNames []string, err error) {
	err = config.Init(configPath)
	if err != nil {
		return
	}

	packNames = config.Vars.Run

	return
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
