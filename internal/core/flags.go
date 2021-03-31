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

var binariesPath string

// ParseFlags ...
func ParseFlags() {
	var flags cliflags.Flags
	flags.NewStringFlag("binaries-path", "custom override for service pack binary location", binariesPathHandler)
	flags.ExecuteHandlers()
}

func binariesPathHandler(v *string) {
	binariesPath = *v // defaults to an empty string, no checks necessary
}

// GetCommands ...
func GetCommands() (cmdSet []*exec.Cmd, err error) {
	// TODO: give any exec errors a familiar format
	workDir, err := os.Getwd()
	if err != nil {
		return
	}
	configPath := filepath.Join(workDir, "config.yml")

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
		// cmd.Args = append(cmd.Args, "--tags='@k-gen'")
		cmd.Args = append(cmd.Args, "--loglevel=DEBUG")
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

func packBinary(name string) (binaryName string, err error) {
	name = strings.ToLower(name)
	if runtime.GOOS == "windows" && !strings.HasSuffix(name, ".exe") {
		name = fmt.Sprintf("%s.exe", name)
	}
	if binariesPath == "" {
		binariesPath = filepath.Join(userHomeDir(), "probr", "binaries") // TODO Load from config.
	}
	binariesPath = strings.Replace(binariesPath, "~", userHomeDir(), 1)
	plugins, _ := hcplugin.Discover(name, binariesPath)
	if len(plugins) != 1 {
		err = fmt.Errorf("Please ensure requested plugin '%s' has been installed to '%s'", name, binariesPath)
		return
	}
	binaryName = plugins[0]

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
