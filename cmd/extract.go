package cmd

import (
	"github.com/spf13/cobra"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract information from a Healthbot Installation",
	Long:  `Extract dsl from Healthbot for e.g. Devices, Device Groups, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.PersistentFlags().StringP("output_directory", "o", ".", "directory where the configuration will be stored")
}
