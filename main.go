package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	plugin "github.com/citihub/probr-sdk/plugin"
	hclog "github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
)

var handshakeConfig_spKubernetes = hcplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "probr.servicepack.probr",
}

// pluginKeyword_spKubernetes is used to identify plugin in pluginmap, and to dispense
var pluginKeyword_spKubernetes = "spProbr"

// pluginMap_spKubernetes is the map of plugins we can dispense.
var pluginMap_spKubernetes = map[string]hcplugin.Plugin{
	pluginKeyword_spKubernetes: &plugin.ServicePackPlugin{},
}

// Location for plugin binaries
var pluginPath_spKubernetes = "./servicepacks/kubernetes/kubernetes"

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "probr",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// We're a host! Start by launching the plugin process.
	cmd := exec.Command(pluginPath_spKubernetes)
	cmd.Args = append(cmd.Args, "--varsfile=config.yml")
	cmd.Args = append(cmd.Args, "--tags=@k-gen")
	client_spKubernetes := hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: handshakeConfig_spKubernetes,
		Plugins:         pluginMap_spKubernetes,
		Cmd:             cmd,
		Logger:          logger,
	})
	defer client_spKubernetes.Kill()

	// Connect via RPC
	rpcClient_spKubernetes, err := client_spKubernetes.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw_spKubernetes, err := rpcClient_spKubernetes.Dispense(pluginKeyword_spKubernetes)
	if err != nil {
		log.Fatal(err)
	}

	// We should have a ServicePack now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	spKubernetes := raw_spKubernetes.(plugin.ServicePack)
	fmt.Println(spKubernetes.Greet())
}
