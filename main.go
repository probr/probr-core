package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/citihub/probr-sdk/plugin"
	hclog "github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig_spKubernetes = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "probr.servicepack.kubernetes",
}
var handshakeConfig_spAPIM = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "probr.servicepack.apim",
}
var handshakeConfig_spProbr = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "probr.servicepack.probr",
}

// pluginKeyword is used to identify plugin in pluginmap, and to dispense
var pluginKeyword_spKubernetes = "spKubernetes"
var pluginKeyword_spAPIM = "spAPIM"
var pluginKeyword_spProbr = "spProbr"

// pluginMap is the map of plugins we can dispense.
var pluginMap_spKubernetes = map[string]plugin.Plugin{
	pluginKeyword_spKubernetes: &probrsdk.ServicePackPlugin{},
}
var pluginMap_spAPIM = map[string]plugin.Plugin{
	pluginKeyword_spAPIM: &probrsdk.ServicePackPlugin{},
}
var pluginMap_spProbr = map[string]plugin.Plugin{
	pluginKeyword_spProbr: &probrsdk.ServicePackPlugin{},
}

// Location for plugin binaries
var pluginPath_spKubernetes = "./servicepacks/kubernetes/kubernetes"
var pluginPath_spAPIM = "./servicepacks/apim/apim"
var pluginPath_spProbr = "./cmd/probr"

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "probr",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// TODO: Load all plugins. Get a collection of probrsdk.ServicePack interface.
	// 	Defer sp.kill
	//  Execute sp.Greet() for all

	// Kubernetes *********************************************************
	// We're a host! Start by launching the plugin process.
	client_spKubernetes := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig_spKubernetes,
		Plugins:         pluginMap_spKubernetes,
		Cmd:             exec.Command(pluginPath_spKubernetes),
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
	spKubernetes := raw_spKubernetes.(probrsdk.ServicePack)
	fmt.Println(spKubernetes.Greet())
	// Kubernetes *********************************************************

	// APIM ***************************************************************
	// We're a host! Start by launching the plugin process.
	client_spAPIM := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig_spAPIM,
		Plugins:         pluginMap_spAPIM,
		Cmd:             exec.Command(pluginPath_spAPIM),
		Logger:          logger,
	})
	defer client_spAPIM.Kill()

	// Connect via RPC
	rpcClient_spAPIM, err := client_spAPIM.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw_spAPIM, err := rpcClient_spAPIM.Dispense(pluginKeyword_spAPIM)
	if err != nil {
		log.Fatal(err)
	}

	// We should have a ServicePack now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	spAPIM := raw_spAPIM.(probrsdk.ServicePack)
	fmt.Println(spAPIM.Greet())
	// APIM ***************************************************************

	// Storage*************************************************************
	// TODO
	// Storage*************************************************************

	// Probr **************************************************************
	// We're a host! Start by launching the plugin process.
	cmd := exec.Command(pluginPath_spProbr)
	//cmd.SysProcAttr = &syscall.SysProcAttr{}
	//cmd.SysProcAttr.CmdLine = fmt.Sprintf("%s %s", pluginPath_spProbr, pluginArgs_spProbr)
	cmd.Args = append(cmd.Args, "--varsfile=config.yml")
	cmd.Args = append(cmd.Args, "--tags=@k-gen")
	client_spProbr := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig_spProbr,
		Plugins:         pluginMap_spProbr,
		Cmd:             cmd,
		Logger:          logger,
	})
	defer client_spProbr.Kill()

	// Connect via RPC
	rpcClient_spProbr, err := client_spProbr.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw_spProbr, err := rpcClient_spProbr.Dispense(pluginKeyword_spProbr)
	if err != nil {
		log.Fatal(err)
	}

	// We should have a ServicePack now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	spProbr := raw_spProbr.(probrsdk.ServicePack)
	fmt.Println(spProbr.Greet())
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter your city: ")
	// city, _ := reader.ReadString('\n')
	// fmt.Print("You live in " + city)
	// fmt.Println(spProbr.Greet())
	// Probr **************************************************************
}
