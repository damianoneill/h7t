package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

func removeBracketsFromJSONObject(src []byte) (bare []byte, err error) {
	// remove all whitespace / carriage returns, etc.
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, src)
	if err != nil {
		return
	}
	bare = bytes.TrimPrefix(buffer.Bytes(), []byte("{"))
	bare = bytes.TrimSuffix(bare, []byte("}"))
	return
}

func appendResourceBody(backup *[]byte, responseBody []byte) (err error) {
	b, err := removeBracketsFromJSONObject(responseBody)
	// append to main json object and add comma in advance of next object being added
	*backup = append(*backup, b...)
	*backup = append(*backup, ',')
	return
}

func createHelper(rc *resty.Client, resource, filename string) (err error) {
	resp, restErr := rc.R().
		SetBasicAuth(ci.Username, ci.Password).
		Get("https://" + ci.Authority + resource)
	if restErr != nil {
		return restErr
	}
	// write it to file
	err = afero.WriteFile(AppFs, filename, resp.Body(), 0777)
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stdout, "Extracted helper files to %v\n", filename)
	return
}

func createBackup(rc *resty.Client, resources []string, filename string) (err error) {
	b := []byte{'{'}
	for _, resource := range resources {
		resp, restErr := rc.R().
			SetBasicAuth(ci.Username, ci.Password).
			Get("https://" + ci.Authority + resource)
		if restErr != nil {
			return restErr
		}
		// ensure we don't get an empty body
		if resp.Body() != nil && resp.StatusCode() == 200 && len(resp.Body()) > 4 {
			err = appendResourceBody(&b, resp.Body())
			if err != nil {
				return
			}
			fmt.Fprintf(os.Stdout, "Including: %v \n", resource)
		} else {
			fmt.Fprintf(os.Stdout, "Excluding (no content): %v \n", resource)
		}
	}

	// need to remove the last comma we added, to make it valid json or the indenter will barf
	// and close the outer object
	b = bytes.TrimSuffix(b, []byte(","))
	b = append(b, '}')

	// pretty print the output
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, b, "", "  ")
	if err != nil {
		return
	}

	// write it to file
	err = afero.WriteFile(AppFs, filename, prettyJSON.Bytes(), 0777)
	if err != nil {
		return
	}
	fmt.Fprintf(os.Stdout, "Extracted config to %v\n", filename)
	return err
}

// extractInstallationCmd represents the extract installation command
var extractInstallationCmd = &cobra.Command{
	Use:   "installation",
	Short: "Extract installation configuration",
	Long:  `Collect from a Healthbot installation the complete configuration and generate a backup file.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jsonResources := []string{
			"/api/v1/devices/",
			"/api/v1/topics/",
			"/api/v1/playbooks/",
			"/api/v1/device-groups/",
			"/api/v1/network-groups/",
			"/api/v1/notifications/",
			"/api/v1/retention-policies/",
			"/api/v1/system-settings/report-generation/destinations/",
			"/api/v1/system-settings/report-generation/reports/",
			"/api/v1/system-settings/schedulers",
			"/api/v1/data-store/grafana/",
			// TODO - extra resources in 2.1?
		}

		t := time.Now()
		err = createBackup(resty.DefaultClient,
			jsonResources,
			cmd.Flag("output_directory").Value.String()+filePathSeperator+"healthbot_backup-"+t.Format("20060102150405")+".json")
		if err != nil {
			return
		}

		// TODO Helper Files
		err = createHelper(resty.DefaultClient,
			"/api/v1/files/helper-files/backup",
			cmd.Flag("output_directory").Value.String()+filePathSeperator+"healthbot_helpers-"+t.Format("20060102150405")+".tar.gz")
		return
	},
}

func init() {
	extractCmd.AddCommand(extractInstallationCmd)
}
