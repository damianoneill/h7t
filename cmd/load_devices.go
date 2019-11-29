package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/damianoneill/h7t/dsl"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

func getDirectoryContents(inputDirectory string) (matches []string, err error) {
	matches, err = afero.Glob(AppFs, inputDirectory+filePathSeperator+"*")
	if err != nil {
		return
	}
	tmp := matches[:0]
	for _, match := range matches {
		isDir, direrr := afero.IsDir(AppFs, match)
		if direrr != nil {
			return
		}
		if !isDir {
			tmp = append(tmp, match)
		}
	}
	return tmp, err
}

func loadThings(thing dsl.Thing, files []string, shouldErase string) (err error) {
	for _, filename := range files {
		err = dsl.ReadThingFromFile(thing, filename, ioutil.ReadFile)
		if err != nil {
			return
		}
		if shouldErase == "true" {
			for _, t := range thing.InnerThings() {
				err = dsl.DeleteThingToResource(resty.DefaultClient, t, ci, true)
				if err != nil {
					return
				}
			}
		} else {
			err = dsl.PostThingToResource(resty.DefaultClient, thing, ci, true)
			if err != nil {
				return
			}
		}
		fmt.Fprintf(os.Stdout, "Updated %v Things from %v to %v \n", thing.Count(), filename, ci.Authority)
	}
	return
}

// loadDevicesCmd represents the devices command
var loadDevicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Load Device configuration",
	Long:  `Load into a Healthbot installation the configuration for the Devices.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		files, err := getDirectoryContents(cmd.Flag("input_directory").Value.String())
		if err != nil {
			return
		}
		err = loadThings(&dsl.Devices{}, files, cmd.Flag("erase").Value.String())
		return
	},
}

func init() {
	loadDevicesCmd.PersistentFlags().BoolP("erase", "e", false, "erase the thing(s) identified in configuration")
	loadCmd.AddCommand(loadDevicesCmd)
}
