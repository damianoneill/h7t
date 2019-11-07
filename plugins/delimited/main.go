package main

import (
	"os"

	"github.com/damianoneill/h7t/dsl"
	"github.com/damianoneill/h7t/plugins"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// DelimitedDevices - is a implementation of Transformer
type DelimitedDevices struct {
	logger hclog.Logger
}

// Devices - returns a list of dsl Devices
func (g *DelimitedDevices) Devices() dsl.Devices {
	g.logger.Debug("message from DelimitedDevices.Devices")
	return dsl.Devices{
		Device: []dsl.Device{dsl.Device{DeviceID: "10.0.0.1"}},
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "delimited",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	transformer := &DelimitedDevices{
		logger: logger,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"transformer": &plugins.TransformerPlugin{Impl: transformer},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
