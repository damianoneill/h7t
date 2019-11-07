package cmd

import (
	"fmt"
	"os"

	"github.com/damianoneill/h7t/dsl"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// extractDevicesCmd represents the devices command
var extractDevicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Extract Device configuration",
	Long:  `Collect from a Healthbot installation the configuration for the Devices.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		devices := dsl.Devices{}
		err = dsl.ExtractThingFromResource(resty.DefaultClient, &devices, ci)
		if err != nil {
			return
		}
		if len(devices.Device) == 0 {
			fmt.Fprintf(os.Stdout, "Zero Devices retrieved, not writing to file \n")
			return
		}
		namedfile := cmd.Flag("output_directory").Value.String() + filePathSeperator + "devices.yml"
		f, err := os.Create(namedfile)
		if err != nil {
			return
		}
		defer f.Close()
		err = dsl.WriteThingToFile(&devices, f)
		if err != nil {
			return
		}
		fmt.Fprintf(os.Stdout, "Wrote %v Devices to %v \n", len(devices.Device), namedfile)
		return f.Sync() // https://www.joeshaw.org/dont-defer-close-on-writable-files/#update-2
	},
}

func init() {
	extractCmd.AddCommand(extractDevicesCmd)
}
