package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/plugin"
	hcplugin "github.com/hashicorp/go-plugin"

	"github.com/citihub/probr-core/internal/core"
)

var packName, varsFile string

func main() {

	cmdSet := parseFlags()

	// Launch the plugin process
	client := core.NewClient(cmdSet[0])
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	rawSP, err := rpcClient.Dispense(plugin.ServicePackPluginName)
	if err != nil {
		log.Fatal(err)
	}

	// We should have a ServicePack now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	servicePack := rawSP.(plugin.ServicePack)
	fmt.Println(servicePack.Greet())
}

func userHomeDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return user.HomeDir
}

func packBinary(name string) string {
	if runtime.GOOS == "windows" {
		name = fmt.Sprintf("%s.exe", name)
	}
	binaryPath := filepath.Join(userHomeDir(), "probr", "binaries")
	plugins, _ := hcplugin.Discover(name, binaryPath)
	if len(plugins) != 1 {
		panic(fmt.Sprintf("Please ensure requested plugin '%s' has been installed to '%s'", name, binaryPath))
	}
	return plugins[0]
}

func parseFlags() []*exec.Cmd {
	var configPath string
	argCount := len(os.Args[1:])
	if argCount < 1 {
		// TODO: deal with this error properly
		log.Fatal("First argument should path to config file")
	} else {
		configPath = os.Args[1]
	}

	packNames := getPackNameFromConfig(configPath)

	var cmdSet []*exec.Cmd
	for _, pack := range packNames {
		cmd := exec.Command(packBinary(pack))
		cmd.Args = append(cmd.Args, fmt.Sprintf("--varsfile=%s", configPath))

		if argCount > 3 {
			// TODO: passing flags to service pack isn't scalable
			cmd.Args = append(cmd.Args, os.Args[3:]...)
		}
		cmdSet = append(cmdSet, cmd)
	}
	return cmdSet
}

func getPackNameFromConfig(configPath string) (packNames []string) {
	config.Init(configPath)
	return config.Vars.Run
}
