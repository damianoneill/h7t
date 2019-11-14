package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/damianoneill/h7t/dsl"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// DeleteHelperFileToResource - Specific op for POSTING to Helper Files
func DeleteHelperFileToResource(rc *resty.Client, filename string, ci dsl.ConnectionInfo) (err error) {
	resp, err := rc.R().
		SetBasicAuth(ci.Username, ci.Password).
		Delete("https://" + ci.Authority + "/api/v1/files/helper-files/" + filepath.Base(filename) + "/")
	if err != nil {
		return
	}
	switch resp.StatusCode() {
	case 204:
		break
	default:
		return errors.New("Problem deleting File: %v " + resp.String())
	}
	fmt.Fprintf(os.Stdout, "Deleted %v to %v \n", filename, ci.Authority)
	return
}

// PostHelperFileToResource - Specific op for POSTING to Helper Files
func PostHelperFileToResource(rc *resty.Client, filename string, ci dsl.ConnectionInfo) (err error) {
	f, err := AppFs.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	resp, err := rc.R().
		SetBasicAuth(ci.Username, ci.Password).
		SetFileReader("up_file", filepath.Base(filename), f). // stripping path from filename
		Post("https://" + ci.Authority + "/api/v1/files/helper-files/" + filepath.Base(filename) + "/")
	if err != nil {
		return
	}
	switch resp.StatusCode() {
	case 200:
		break
	default:
		return errors.New("Problem uploading File: %v " + resp.String())
	}
	fmt.Fprintf(os.Stdout, "Uploaded %v to %v \n", filename, ci.Authority)
	return
}

// helperFilesCmd represents the Helper Files command
var helperFilesCmd = &cobra.Command{
	Use:   "helper-files",
	Short: "Load Helper Files",
	Long:  `Load into a Healthbot installation the Helper Files e.g. any required python files.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		files, err := getDirectoryContents(cmd.Flag("input_directory").Value.String())
		if err != nil {
			return
		}
		for _, filename := range files {
			if cmd.Flag("erase").Value.String() == "true" {
				err = DeleteHelperFileToResource(resty.DefaultClient, filename, ci)
				if err != nil {
					return
				}
			} else {
				err = PostHelperFileToResource(resty.DefaultClient, filename, ci)
				if err != nil {
					return
				}
			}
		}
		return
	},
}

func init() {
	loadCmd.AddCommand(helperFilesCmd)
}
