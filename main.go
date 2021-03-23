package main

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	plugin "github.com/citihub/probr-sdk/plugin"
	hcplugin "github.com/hashicorp/go-plugin"

	"github.com/citihub/probr-core/internal/core"
)

var cmd string

func main() {
	packName := "azureapim"
	cmd := exec.Command(packBinary(packName))
	cmd.Args = append(cmd.Args, "--varsfile=config.yml")
	cmd.Args = append(cmd.Args, "--tags=@k-gen")

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
