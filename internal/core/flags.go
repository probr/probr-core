package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/citihub/probr-sdk/config"
	hcplugin "github.com/hashicorp/go-plugin"

	cliflags "github.com/citihub/probr-sdk/cli_flags"
)

// BinariesPath represents the path where service pack binaries are installed
var BinariesPath string

// ParseFlags ...
func ParseFlags() {
	var flags cliflags.Flags
	flags.NewStringFlag("binaries-path", "Location for service pack binaries. If not provided, default value is: [UserHomeDir]/probr/binaries", binariesPathHandler)
	flags.ExecuteHandlers()
}

func binariesPathHandler(v *string) {
	BinariesPath = *v // defaults to an empty string, no checks necessary
}

// GetCommands ...
func GetCommands() (cmdSet []*exec.Cmd, err error) {
	// TODO: give any exec errors a familiar format
	configPath, err := getConfigPath()
	if err != nil {
		return
	}
	packNames, err := GetPackNameFromConfig()
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
	return
}

func userHomeDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return user.HomeDir
}

// GetPackBinary finds provided service pack in installation folder and return binary name
func GetPackBinary(name string) (binaryName string, err error) {
	name = strings.ToLower(name)
	if runtime.GOOS == "windows" && !strings.HasSuffix(name, ".exe") {
		name = fmt.Sprintf("%s.exe", name)
	}
	if BinariesPath == "" {
		BinariesPath = filepath.Join(userHomeDir(), "probr", "binaries") // TODO Load from config.
	}
	BinariesPath = strings.Replace(BinariesPath, "~", userHomeDir(), 1)
	plugins, _ := hcplugin.Discover(name, BinariesPath)
	if len(plugins) != 1 {
		err = fmt.Errorf("Please ensure requested plugin '%s' has been installed to '%s'", name, BinariesPath)
		return
	}
	binaryName = plugins[0]

	return
}

func getConfigPath() (configPath string, err error) {
	workDir, err := os.Getwd()
	if err != nil {
		return
	}
	configPath = filepath.Join(workDir, "config.yml")

	return
}

// GetPackNameFromConfig returns all service packs declared in config file
func GetPackNameFromConfig() (packNames []string, err error) {
	configPath, err := getConfigPath()
	if err != nil {
		return
	}

	err = config.Init(configPath)
	if err != nil {
		return
	}

	packNames = config.Vars.Run

	return
}
