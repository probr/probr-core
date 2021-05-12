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

	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/plugin"
	"github.com/citihub/probr-sdk/probeengine"
	"github.com/citihub/probr-sdk/utils"

	hclog "github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
)

// BinariesPath represents the path where service pack binaries are installed
// Must be a pointer to accept the flag when it is set
var BinariesPath *string

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
	configPath, err := getConfigPath()
	if err != nil {
		return
	}
	packNames, err := GetPackNames()
	if err != nil {
		return
	}

	for _, pack := range packNames {
		binaryName, binErr := GetPackBinary(pack)
		if binErr != nil {
			err = binErr
			break
		}
		cmd := exec.Command(binaryName)
		cmd.Args = append(cmd.Args, fmt.Sprintf("--varsfile=%s", configPath))
		cmdSet = append(cmdSet, cmd)
	}
	if err == nil && len(cmdSet) == 0 {
		err = utils.ReformatError("No valid service packs specified")
	}
	return
}

// UserHomeDir provides the OS-aware user home directory
// TODO: move this to SDK
func UserHomeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return dirname
}

// GetPackBinary finds provided service pack in installation folder and return binary name
func GetPackBinary(name string) (binaryName string, err error) {
	name = filepath.Base(strings.ToLower(name)) // in some cases a filepath may arrive here instead of the base name
	if runtime.GOOS == "windows" && !strings.HasSuffix(name, ".exe") {
		name = fmt.Sprintf("%s.exe", name)
	}
	*BinariesPath = strings.Replace(*BinariesPath, "~", UserHomeDir(), 1)
	plugins, _ := hcplugin.Discover(name, *BinariesPath)
	if len(plugins) != 1 {
		err = fmt.Errorf("failed to locate requested plugin '%s'", name)
		return
	}
	binaryName = plugins[0]

	return
}

func getConfigPath() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(workDir, "config.yml"), nil
}

// GetPackNames returns all service packs declared in config file
func GetPackNames() (packNames []string, err error) {
	if AllPacks != nil && *AllPacks {
		return hcplugin.Discover("*", *BinariesPath)
	}
	packNames, err = getPackNamesFromConfig()
	return
}

func getPackNamesFromConfig() ([]string, error) {
	type simpleVars struct {
		Run []string `yaml:"Run"`
	}
	var vars simpleVars

	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	configDecoder, file, err := config.NewConfigDecoder(configPath)
	if err != nil {
		return nil, err
	}

	err = configDecoder.Decode(&vars)
	file.Close()
	return vars.Run, err

}
