package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/damianoneill/h7t/dsl"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// summariseInstallationCmd represents the installation command
var summariseInstallationCmd = &cobra.Command{
	Use:   "installation",
	Short: "Summarise information collected from a Healthbot installation",
	Long:  `Generates counts from different Healthbot dsl things for e.g. Devices, Device Groups etc.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return summariseInstallation(ci)
	},
}

// NewTable - provides a blank table for rendering.
func NewTable(out io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(out)
	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	return table
}

func collectSystemDetails(rc *resty.Client, ci dsl.ConnectionInfo, stdout io.Writer) (err error) {
	sd := dsl.SystemDetails{}
	err = dsl.ExtractThingFromResource(rc, &sd, ci)
	if err != nil {
		return
	}
	fmt.Fprintln(stdout, "")
	fmt.Fprintf(stdout, "Healthbot Authority: %s \n", ci.Authority)
	fmt.Fprintf(stdout, "Healthbot Version: %s \n", sd.Version)
	fmt.Fprintf(stdout, "Healthbot Time: %s \n", sd.ServerTime)
	return
}

func collectDeviceFacts(rc *resty.Client, ci dsl.ConnectionInfo, stdout io.Writer) (df dsl.DeviceFacts, err error) {
	err = dsl.ExtractThingFromResource(rc, &df, ci)
	if err != nil {
		return
	}
	fmt.Fprintln(stdout, "")
	fmt.Fprintf(stdout, "No of Managed Devices: %v \n", len(df))
	fmt.Fprintln(stdout, "")
	return
}

func renderDeviceTable(w io.Writer, df dsl.DeviceFacts) {
	table := NewTable(w)
	table.SetHeader([]string{"Device Id", "Platform", "Release", "Serial Number"})
	table.Append([]string{"", "", "", ""})
	for _, fact := range df {
		table.Append([]string{fact.DeviceID, fact.Facts.Platform, fact.Facts.Release, fact.Facts.SerialNumber})
	}
	table.Render() // Send output
}

func summariseInstallation(ci dsl.ConnectionInfo) (err error) {

	err = collectSystemDetails(resty.DefaultClient, ci, os.Stdout)
	if err != nil {
		return
	}

	df, err := collectDeviceFacts(resty.DefaultClient, ci, os.Stdout)
	if err != nil {
		return
	}

	renderDeviceTable(os.Stdout, df)

	return
}

func init() {
	summariseCmd.AddCommand(summariseInstallationCmd)
}
