package cmd

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
)

var logLevel = hclog.Info

// transformCmd represents the transform command
var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "Transform things from proprietary formats into Healthbot dsl format",
	Long:  `Transform customer content into dsl things for e.g. devices.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("verbose").Value.String() == "true" {
			logLevel = hclog.Debug
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(transformCmd)
	transformCmd.PersistentFlags().StringP("output_directory", "o", ".", "directory where the configuration will be written to")
	transformCmd.PersistentFlags().StringP("input_directory", "i", ".", "directory where the configuration will be loaded from")
	transformCmd.PersistentFlags().String("plugin", "csv", "name of the plugin to be used")
}
