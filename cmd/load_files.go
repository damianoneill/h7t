package cmd

import (
	"github.com/spf13/cobra"
)

// loadFilesCmd represents the files command
var loadFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Load file formats into a Healthbot Installation",
	Long: `Load Files into Healthbot for e.g. Helper, Playbook or Rule Files etc..

Files sub-commands work by iterating over all files in the input directory, it will filter based on the file type.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	loadCmd.AddCommand(loadFilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
