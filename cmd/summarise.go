package cmd

import (
	"github.com/spf13/cobra"
)

// summariseCmd represents the summarise command
var summariseCmd = &cobra.Command{
	Use:   "summarise",
	Short: "Summarise information from a Healthbot Installation",
	Long:  `Summarise dsl from Healthbot for e.g. Devices, Device Groups, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(summariseCmd)
}
