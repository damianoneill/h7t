package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// loadFilesPlaybooksCmd represents the playbook command
var loadFilesPlaybooksCmd = &cobra.Command{
	Use:   "playbooks",
	Short: "Load Playbook Files",
	Long: `Load into a Healthbot installation the Playbook Files.
	
The erase option is not supported for this command as the Playbook cannot be indentified from the file name.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		files, err := getDirectoryContents(cmd.Flag("input_directory").Value.String())
		if err != nil {
			return
		}
		for _, filename := range files {
			// ignore everything but playbooks
			if strings.HasSuffix(filename, ".playbook") {
				err = PostFileToResource(resty.DefaultClient, filename, "/api/v1/playbooks/", "playbooks", ci, true)
				if err != nil {
					return
				}
			}
		}
		return
	},
}

func init() {
	loadFilesCmd.AddCommand(loadFilesPlaybooksCmd)
}
