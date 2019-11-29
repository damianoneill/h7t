package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// loadFilesRulesCmd represents the rules command
var loadFilesRulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Load Rule Files",
	Long: `Load into a Healthbot installation the Rule Files.
	
The erase option is not supported for this command as the Rule cannot be indentified from the file name.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		files, err := getDirectoryContents(cmd.Flag("input_directory").Value.String())
		if err != nil {
			return
		}
		for _, filename := range files {
			// ignore everything but rules
			if strings.HasSuffix(filename, ".rule") {
				err = PostFileToResource(resty.DefaultClient, filename, "/api/v1/topics/", "topics", ci, true)
				if err != nil {
					return
				}
			}
		}
		return
	},
}

func init() {
	loadFilesCmd.AddCommand(loadFilesRulesCmd)
}
