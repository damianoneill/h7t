package cmd

import (
	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load information into a Healthbot Installation",
	Long: `Load dsl into Healthbot for e.g. Devices, Device Groups, etc..

Load sub-commands work by iterating over all files in the input directory`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.PersistentFlags().StringP("input_directory", "i", ".", "directory where the configuration will be loaded from")
	loadCmd.PersistentFlags().BoolP("erase", "e", false, "erase the thing(s) identified in configuration")
}
