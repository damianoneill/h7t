package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jszwec/csvutil"

	"github.com/damianoneill/h7t/dsl"
	"github.com/damianoneill/h7t/plugins"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// CsvDevices - is a implementation of Transformer
type CsvDevices struct {
	logger hclog.Logger
}

// ReadCsvFile - using csvutil to unmarshal, it expects a slice of types
func ReadCsvFile(filename string, readfile func(filename string) ([]byte, error)) (devices []dsl.Device, err error) {
	b, err := readfile(filename)
	if err != nil {
		return
	}
	err = csvutil.Unmarshal(b, &devices)
	return
}

// Devices - returns a list of dsl Devices
func (g *CsvDevices) Devices(args plugins.Arguments) (devices dsl.Devices, err error) {
	g.logger.Debug("args:", "inputDirectory", args.InputDirectory)

	// can be used by plugin providers to handle arguments
	g.logger.Debug("args:", "command line arguments", args.CmdLineArgs)

	// get all files in the input directory
	files, err := filepath.Glob(args.InputDirectory + string(filepath.Separator) + "*")
	if err != nil {
		return
	}

	for _, filename := range files {
		var d []dsl.Device
		d, err = ReadCsvFile(filename, ioutil.ReadFile)
		if err != nil {
			return
		}
		devices.Device = append(devices.Device, d...)
	}

	return
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "csv",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	transformer := &CsvDevices{
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
