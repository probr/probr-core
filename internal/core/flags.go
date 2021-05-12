package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	hcplugin "github.com/hashicorp/go-plugin"

	"github.com/citihub/probr-sdk/config"
)

// BinariesPath represents the path where service pack binaries are installed
// Must be a pointer to accept the flag when it is set
var BinariesPath *string

// Verbose is a CLI option to increase output detail
var Verbose *bool

// AllPacks is a CLI option to target all installed packs, instead of just those specified in config.yml
var AllPacks *bool

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
	binaryName = filepath.Base(plugins[0])

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

	type simpleVars struct {
		Run []string `yaml:"Run"`
	}
	var vars simpleVars

	configPath, err := getConfigPath()
	if err != nil {
		return
	}

	configDecoder, file, err := config.NewConfigDecoder(configPath)
	if err != nil {
		return
	}

	err = configDecoder.Decode(&vars)
	file.Close()
	packNames = vars.Run
	return
}

// TODO: Seems like these functions could use with some topical reorganization
