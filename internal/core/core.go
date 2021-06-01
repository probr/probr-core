package core

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

	"github.com/probr/probr-sdk/plugin"
	"github.com/probr/probr-sdk/probeengine"
	"github.com/probr/probr-sdk/utils"
	"github.com/probr/probr/internal/config"

	hclog "github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
)

// Verbose is a CLI option to increase output detail
var Verbose *bool

// AllPacks is a CLI option to target all installed packs, instead of just those specified in config.yml
var AllPacks *bool

// NewClient client handles the lifecycle of a plugin application
// Plugin hosts should use one Client for each plugin executable
// (this is different from the client that manages gRPC)
func NewClient(cmd *exec.Cmd) *hcplugin.Client {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   plugin.ServicePackPluginName,
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	var pluginMap = map[string]hcplugin.Plugin{
		plugin.ServicePackPluginName: &plugin.ServicePackPlugin{},
	}
	var handshakeConfig = plugin.GetHandshakeConfig()
	return hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		Logger:          logger,
	})
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
// Ref: https://golangcode.com/handle-ctrl-c-exit-in-terminal/
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("Execution aborted - %v", "SIGTERM")
		probeengine.CleanupTmp()
		os.Exit(0)
	}()
}

// GetCommands ...
func GetCommands() (cmdSet []*exec.Cmd, err error) {
	// TODO: give any exec errors a familiar format

	for _, pack := range config.Vars.Run {
		binaryName, binErr := GetPackBinary(pack)
		if binErr != nil {
			err = binErr
			break
		}
		cmd := exec.Command(binaryName)
		cmd.Args = append(cmd.Args, fmt.Sprintf("--varsfile=%s", *config.Vars.VarsFile))
		cmdSet = append(cmdSet, cmd)
	}
	log.Printf("BIN: %s", config.Vars.BinariesPath)
	if err == nil && len(cmdSet) == 0 {
		available, _ := hcplugin.Discover("*", config.Vars.BinariesPath)
		err = utils.ReformatError("No valid service packs specified. Requested: %v, Available: %v", config.Vars.Run, available)
	}
	return
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

// GetPackNames returns all service packs declared in config file
func GetPackNames() (packNames []string, err error) {
	if err != nil || (AllPacks != nil && *AllPacks) {
		return hcplugin.Discover("*", config.Vars.BinariesPath)
	}
	return config.Vars.Run, nil
}
