package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	plugin "github.com/citihub/probr-sdk/plugin"
	hcplugin "github.com/hashicorp/go-plugin"

	"github.com/citihub/probr-core/internal/core"
)

var packName, varsFile string

func main() {
	argCount := len(os.Args[1:])
	if argCount < 2 {
		// TODO: deal with this error properly
		log.Fatal("First argument should be the name of the service pack, second should be path to config file")
	} else {
		packName = os.Args[1]
		varsFile = os.Args[2]
	}
	cmd := exec.Command(packBinary(packName))
	cmd.Args = append(cmd.Args, fmt.Sprintf("--varsfile=%s", varsFile))

	if argCount > 3 {
		cmd.Args = append(cmd.Args, os.Args[3:]...)
	}

	// Launch the plugin process
	client := core.NewClient(cmd, packName)
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	rawSP, err := rpcClient.Dispense(packName)
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
