package cmd

import (
	"github.com/damianoneill/h7t/dsl"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// extractDeviceGroupsCmd represents the extract DeviceGroups command
var extractDeviceGroupsCmd = &cobra.Command{
	Use:   "device-groups",
	Short: "Extract Device Groups configuration",
	Long:  `Collect from a Healthbot installation the configuration for the Devices Groups.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		dg := dsl.DeviceGroups{}
		err = dsl.ExtractThingFromResource(resty.DefaultClient, &dg, ci)
		if err != nil {
			return
		}
		return WriteThingsToFile(&dg, cmd.Flag("output_directory").Value.String()+filePathSeperator+"device-groups.yml")
	},
}

func init() {
	extractCmd.AddCommand(extractDeviceGroupsCmd)
}
