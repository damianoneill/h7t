package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/damianoneill/h7t/dsl"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Load Device configuration",
	Long:  `Load into a Healthbot installation the configuration for the Devices.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		inputDirectory := cmd.Flag("input_directory").Value.String()
		shouldErase := cmd.Flag("erase").Value.String()
		files, err := filepath.Glob(inputDirectory + filePathSeperator + "*")
		if err != nil {
			return
		}
		for _, filename := range files {
			devices := dsl.Devices{}
			err = dsl.ReadThingFromFile(&devices, filename, ioutil.ReadFile)
			if err != nil {
				return
			}

			if shouldErase == "true" {
				for _, device := range devices.Device {
					err = dsl.DeleteThingToResource(resty.DefaultClient, &device, ci, true)
					if err != nil {
						return
					}
				}
			} else {
				err = dsl.PostThingToResource(resty.DefaultClient, &devices, ci, true)
				if err != nil {
					return
				}
			}

			fmt.Fprintf(os.Stdout, "Updated %v Devices from %v to %v \n", len(devices.Device), filename, ci.Authority)
		}
		return
	},
}

func init() {
	loadCmd.AddCommand(devicesCmd)
}
