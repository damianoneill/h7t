package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/damianoneill/h7t/pkg/dsl"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// WriteThingsToFile - common function used by commands creating yaml things
func WriteThingsToFile(thing dsl.Thing, namedfile string) (err error) {
	if thing.Count() == 0 {
		return errors.New("zero Things, not writing to file")
	}
	f, err := AppFs.Create(namedfile)
	if err != nil {
		return
	}
	defer f.Close()
	err = dsl.WriteThingToFile(thing, f)
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stdout, "Wrote %v Things to %v \n", thing.Count(), namedfile)
	return f.Sync() // https://www.joeshaw.org/dont-defer-close-on-writable-files/#update-2
}

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
		return WriteThingsToFile(&devices, cmd.Flag("output_directory").Value.String()+filePathSeperator+"devices.yml")
	},
}

func init() {
	extractCmd.AddCommand(extractDevicesCmd)
}
