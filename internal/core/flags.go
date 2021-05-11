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

	hcplugin "github.com/hashicorp/go-plugin"

	"github.com/citihub/probr-sdk/config"
)

// BinariesPath represents the path where service pack binaries are installed
// Must be a pointer to accept the flag when it is set
var BinariesPath *string

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
	if *BinariesPath == "" {
		*BinariesPath = filepath.Join(userHomeDir(), "probr", "binaries") // TODO Load from config.
	}
	*BinariesPath = strings.Replace(*BinariesPath, "~", userHomeDir(), 1)
	plugins, _ := hcplugin.Discover(name, *BinariesPath)
	if len(plugins) != 1 {
		err = fmt.Errorf("Please ensure requested plugin '%s' has been installed to '%s'", name, *BinariesPath)
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

// GetPackNameFromConfig returns all service packs declared in config file
func GetPackNameFromConfig() (packNames []string, err error) {
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
	log.Printf("Found packs %v in config file: %s", packNames, configPath)
	return
}

// TODO: Seems like these functions could use with some topical reorganization
