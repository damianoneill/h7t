package cmd

import (
	"fmt"
	"github.com/damianoneill/h7t/dsl"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Load Device configuration",
	Long:  `Load into a Healthbot installation the configuration for the Devices.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		inputDirectory := cmd.Flag("input_directory").Value.String()
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

			err = dsl.PostThingToResource(resty.DefaultClient, &devices, ci)
			if err != nil {
				return
			}
			fmt.Fprintf(os.Stdout, "Loaded %v Devices from %v to %v \n", len(devices.Device), filename, ci.Authority)
		}
		return
	},
}

func init() {
	loadCmd.AddCommand(devicesCmd)
}
