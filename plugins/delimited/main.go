package main

import (
	"errors"
	"os"
	"regexp"

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
func (g *DelimitedDevices) Devices(args plugins.Arguments) (devices dsl.Devices, err error) {
	g.logger.Debug("args:", "inputDirectory", args.InputDirectory)
	g.logger.Debug("args:", "command line arguments", args.CmdLineArgs)
	if len(args.CmdLineArgs) != 1 {
		err = errors.New("delimited plugin requires regex in quotes \" \" as first argument")
	}
	delimiter := regexp.MustCompile(args.CmdLineArgs[0])
	g.logger.Info("delimiter", "split", delimiter.Split("banana   split", -1))
	devices = dsl.Devices{
		Device: []dsl.Device{dsl.Device{DeviceID: "10.0.0.1"}},
	}
	return
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

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
