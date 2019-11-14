package cmd

import (
	"github.com/damianoneill/h7t/dsl"
	"github.com/spf13/cobra"
)

// loadDeviceGroupsCmd represents the load Device Group command
var loadDeviceGroupsCmd = &cobra.Command{
	Use:   "device-groups",
	Short: "Load Device Groups configuration",
	Long:  `Load into a Healthbot installation the configuration for the Device Groups.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		files, err := getDirectoryContents(cmd.Flag("input_directory").Value.String())
		if err != nil {
			return
		}
		err = loadThings(&dsl.DeviceGroups{}, files, cmd.Flag("erase").Value.String())
		return
	},
}

func init() {
	loadCmd.AddCommand(loadDeviceGroupsCmd)
}
