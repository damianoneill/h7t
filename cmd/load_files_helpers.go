package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/damianoneill/h7t/dsl"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// DeleteFileToResource - Specific op for DELETING Files
func DeleteFileToResource(rc *resty.Client, filename, path string, ci dsl.ConnectionInfo) (err error) {
	t, err := dsl.GetToken(rc, ci)
	if err != nil {
		return err
	}
	resp, err := rc.R().
		SetAuthToken(t.AccessToken).
		SetBasicAuth(ci.Username, ci.Password).
		Delete("https://" + ci.Authority + path)
	if err != nil {
		return
	}
	switch resp.StatusCode() {
	case http.StatusNoContent:
		break
	default:
		return errors.New("Problem deleting File: %v " + resp.String())
	}
	fmt.Fprintf(os.Stdout, "Deleted %v to %v \n", filename, ci.Authority)
	return
}

// PostFileToResource - Specific op for POSTING Files
func PostFileToResource(rc *resty.Client, filename, path, paramName string, ci dsl.ConnectionInfo, shouldCommit bool) (err error) {
	f, err := AppFs.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	t, err := dsl.GetToken(rc, ci)
	if err != nil {
		return err
	}
	resp, err := rc.R().
		SetAuthToken(t.AccessToken).
		SetBasicAuth(ci.Username, ci.Password).
		SetFileReader(paramName, filepath.Base(filename), f). // stripping path from filename
		Post("https://" + ci.Authority + path)
	if err != nil {
		return
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		break
	default:
		return errors.New("Problem uploading File: %v " + resp.String())
	}
	if shouldCommit {
		_, err = rc.R().
			SetAuthToken(t.AccessToken).
			SetBasicAuth(ci.Username, ci.Password).
			Post("https://" + ci.Authority + "/api/v1/configuration/")
		if err != nil {
			return
		}
	}
	fmt.Fprintf(os.Stdout, "Uploaded %v to %v \n", filename, ci.Authority)
	return err
}

// loadFilesHelpersCmd represents the Helper Files command
var loadFilesHelpersCmd = &cobra.Command{
	Use:   "helper-files",
	Short: "Load Helper Files",
	Long: `Load into a Healthbot installation the Helper Files e.g. any required python files.

This will exclude any playbook and rule files in the input directory.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		files, err := getDirectoryContents(cmd.Flag("input_directory").Value.String())
		if err != nil {
			return
		}
		for _, filename := range files {
			// ignore playbook and rule files
			if !strings.HasSuffix(filename, ".playbook") && !strings.HasSuffix(filename, ".rule") {
				if cmd.Flag("erase").Value.String() == TRUE {
					err = DeleteFileToResource(resty.DefaultClient, filename, "/api/v1/files/helper-files/"+filepath.Base(filename)+"/", ci)
					if err != nil {
						return
					}
				} else {
					err = PostFileToResource(resty.DefaultClient,
						filename, "/api/v1/files/helper-files/"+filepath.Base(filename)+"/", "up_file", ci, false)
					if err != nil {
						return
					}
				}
			}
		}
		return
	},
}

func init() {
	loadFilesHelpersCmd.PersistentFlags().BoolP("erase", "e", false, "erase the thing(s) identified in configuration")
	loadFilesCmd.AddCommand(loadFilesHelpersCmd)
}
