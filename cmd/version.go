package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the current build information",
	Long:  "Version, Commit and Date will be output from the Build Info.",
	Run: func(cmd *cobra.Command, args []string) {
		Version(os.Stdout)
	},
}

// Version outputs a formatted version message to the passed writer.
func Version(w io.Writer) {
	fmt.Fprintf(w, "%v, commit %v, built at %v \n", version, commit, date)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
