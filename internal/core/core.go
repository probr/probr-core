package core

import (
	"os"
	"os/exec"

	"github.com/citihub/probr-sdk/plugin"

	hclog "github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
)

var logger hclog.Logger

// NewClient client handles the lifecycle of a plugin application
// Plugin hosts should use one Client for each plugin executable
// (this is different from the client that manages gRPC)
func NewClient(cmd *exec.Cmd) *hcplugin.Client {
	// TODO: Move reusable code blocks to SDK, such as: GetPluginMap; GetHandshakeConfig;
	logger = hclog.New(&hclog.LoggerOptions{
		Name:   plugin.ServicePackPluginName,
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	var pluginMap = map[string]hcplugin.Plugin{
		plugin.ServicePackPluginName: &plugin.ServicePackPlugin{},
	}
	var handshakeConfig = hcplugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "PROBR_MAGIC_COOKIE",
		MagicCookieValue: "probr.servicepack",
	}
	return hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		Logger:          logger,
	})
}
