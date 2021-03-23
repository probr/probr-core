package core

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/citihub/probr-sdk/plugin"

	hclog "github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
)

var logger hclog.Logger

// A client handles the lifecycle of a plugin application
// Plugin hosts should use one Client for each plugin executable
// (this is different from the client that manages gRPC)
func NewClient(cmd *exec.Cmd, packName string) *hcplugin.Client {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:   packName,
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	var pluginMap = map[string]hcplugin.Plugin{
		packName: &plugin.ServicePackPlugin{},
	}
	var handshakeConfig = hcplugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: fmt.Sprintf("probr.servicepack.%s", packName),
	}
	return hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		Logger:          logger,
	})
}
