package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/damianoneill/h7t/plugins"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "delimited",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"transformer": &plugins.TransformerPlugin{},
}

// devicesCmd represents the devices command
var tranformDevicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Transform Devices configuration",
	Long:  `Transform Devices configurations from a proprietary format into the dsl format using a bundled plugin.`,
	Run: func(cmd *cobra.Command, args []string) {
		// inputDirectory := cmd.Flag("input_directory").Value.String()
		// outputDirectory := cmd.Flag("output_directory").Value.String()
		p := cmd.Flag("plugin").Value.String()
		fmt.Fprintf(os.Stdout, "Plugin: %v \n", p)

		// Create an hclog.Logger
		logger := hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Output: os.Stdout,
			Level:  logLevel,
		})

		// We're a host! Start by launching the plugin process.
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         pluginMap,
			Cmd:             exec.Command("./plugins/delimited/transformer"),
			Logger:          logger,
		})
		defer client.Kill()

		// Connect via RPC
		rpcClient, err := client.Client()
		if err != nil {
			log.Fatal(err)
		}

		// Request the plugin
		raw, err := rpcClient.Dispense("transformer")
		if err != nil {
			log.Fatal(err)
		}

		transformer := raw.(plugins.Transformer)
		fmt.Fprintf(os.Stdout, "Devices: %v \n", transformer.Devices(args))
	},
}

func init() {
	transformCmd.AddCommand(tranformDevicesCmd)
}
