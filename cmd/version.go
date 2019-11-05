package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the current build information",
	Long:  "Version, Commit and Date will be output from the Build Info.",
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("%v, commit %v, built at %v \n", bi.version, bi.commit, bi.date)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
