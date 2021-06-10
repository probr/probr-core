package run

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	hcplugin "github.com/hashicorp/go-plugin"
	sdkconfig "github.com/probr/probr-sdk/config"
	"github.com/probr/probr-sdk/logging"
	"github.com/probr/probr-sdk/plugin"
	"github.com/probr/probr-sdk/probeengine"
	"github.com/probr/probr-sdk/utils"

	"github.com/probr/probr/internal/config"
)

// CLIContext executes all plugins with handling for the command line
func CLIContext() {
	// Setup for handling SIGTERM (Ctrl+C)
	setupCloseHandler()

	cmdSet, err := getCommands()
	if err != nil {
		log.Printf("Error loading plugins from config: %s", err)
		os.Exit(2)
	}

	// Run all plugins
	if err := AllPlugins(cmdSet); err != nil {
		log.Printf("[INFO] Output directory: %s", sdkconfig.GlobalConfig.WriteDirectory)
		switch e := err.(type) {
		case *ServicePackErrors:
			log.Printf("[ERROR] %d out of %d test service packs failed. %v", len(e.Errors), len(cmdSet), e)
			os.Exit(1) // At least one service pack failed
		default:
			log.Print(utils.ReformatError(err.Error()))
			os.Exit(2) // Internal error
		}
	}
	log.Printf("[INFO] No errors encountered during plugin execution. Output directory: %s", sdkconfig.GlobalConfig.WriteDirectory)
	os.Exit(0)
}

// AllPlugins executes specified plugins in a loop
func AllPlugins(cmdSet []*exec.Cmd) (err error) {
	spErrors := make([]ServicePackError, 0) // This will store any plugin errors received during execution

	for _, cmd := range cmdSet {
		spErrors, err = Plugin(cmd, spErrors)
		if err != nil {
			return
		}
	}

	if len(spErrors) > 0 {
		// Return all service pack errors to main
		err = &ServicePackErrors{
			Errors: spErrors,
		}
	}
	return
}

// Plugin executes single plugin based on the provided command
func Plugin(cmd *exec.Cmd, spErrors []ServicePackError) ([]ServicePackError, error) {
	// Launch the plugin process
	client := newClient(cmd)
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
		spErr := ServicePackError{
			ServicePack: cmd.String(), // TODO: retrieve service pack name from interface function
			Err:         response,
		}
		spErrors = append(spErrors, spErr)
	} else {
		log.Printf("[INFO] Probes all completed with successful results")
	}
	return spErrors, nil
}

// GetPackBinary finds provided service pack in installation folder and return binary name
func GetPackBinary(name string) (binaryName string, err error) {
	name = filepath.Base(strings.ToLower(name)) // in some cases a filepath may arrive here instead of the base name
	if runtime.GOOS == "windows" && !strings.HasSuffix(name, ".exe") {
		name = fmt.Sprintf("%s.exe", name)
	}
	home, _ := os.UserHomeDir()
	config.Vars.BinariesPath = strings.Replace(config.Vars.BinariesPath, "~", home, 1)
	plugins, _ := hcplugin.Discover(name, config.Vars.BinariesPath)
	if len(plugins) != 1 {
		err = fmt.Errorf("failed to locate requested plugin '%s' at path '%s'", name, config.Vars.BinariesPath)
		return
	}
	binaryName = plugins[0]

	return
}

// setupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
// Ref: https://golangcode.com/handle-ctrl-c-exit-in-terminal/
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("Execution aborted - %v", "SIGTERM")
		probeengine.CleanupTmp()
		os.Exit(0)
	}()
}

func getCommands() (cmdSet []*exec.Cmd, err error) {
	// TODO: give any exec errors a familiar format

	for _, pack := range config.Vars.Run {
		cmd, err := getCommand(pack)
		if err != nil {
			break
		}
		cmdSet = append(cmdSet, cmd)
	}
	log.Printf("[DEBUG] Using bin: %s", config.Vars.BinariesPath)
	if err == nil && len(cmdSet) == 0 {
		available, _ := hcplugin.Discover("*", config.Vars.BinariesPath)
		err = utils.ReformatError("No valid service packs specified. Requested: %v, Available: %v", config.Vars.Run, available)
	}
	return
}

// TODO
func getCommand(pack string) (cmd *exec.Cmd, err error) {
	binaryName, binErr := GetPackBinary(pack)
	if binErr != nil {
		err = binErr
		return
	}
	cmd = exec.Command(binaryName)
	cmd.Args = append(cmd.Args, fmt.Sprintf("--varsfile=%s", *config.Vars.VarsFile))
	return
}

// newClient client handles the lifecycle of a plugin application
// Plugin hosts should use one Client for each plugin executable
// (this is different from the client that manages gRPC)
func newClient(cmd *exec.Cmd) *hcplugin.Client {
	var pluginMap = map[string]hcplugin.Plugin{
		plugin.ServicePackPluginName: &plugin.ServicePackPlugin{},
	}
	var handshakeConfig = plugin.GetHandshakeConfig()
	return hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		Logger:          logging.GetLogger("core"),
	})
}
