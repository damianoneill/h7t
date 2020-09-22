package cmd

import (
	"github.com/spf13/cobra"
)

// summarizeCmd represents the summarize command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarize information from a Healthbot Installation",
	Long:  `Summarize dsl from Healthbot for e.g. Devices, Device Groups, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(summarizeCmd)
}
