package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/damianoneill/h7t/pkg/dsl"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// summarizeInstallationCmd represents the installation command
var summarizeInstallationCmd = &cobra.Command{
	Use:   "installation",
	Short: "Summarize information collected from a Healthbot installation",
	Long:  `Generates counts from different Healthbot dsl things for e.g. Devices, Device Groups etc.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return summarizeInstallation(ci)
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

func collectDeviceGroups(rc *resty.Client, ci dsl.ConnectionInfo, stdout io.Writer) (dg dsl.DeviceGroups, err error) {
	err = dsl.ExtractThingFromResource(rc, &dg, ci)
	if err != nil {
		return
	}
	fmt.Fprintln(stdout, "")
	fmt.Fprintf(stdout, "No of Device Groups: %v \n", len(dg.DeviceGroup))
	fmt.Fprintln(stdout, "")
	return
}

func renderDeviceTable(w io.Writer, df dsl.DeviceFacts) {
	table := NewTable(w)
	table.SetHeader([]string{"Device Id", "Platform", "Release", "Serial Number"})
	table.Append([]string{"", "", "", ""})
	for i := range df {
		table.Append([]string{df[i].DeviceID, df[i].Facts.Platform, df[i].Facts.Release, df[i].Facts.SerialNumber})
	}
	table.Render() // Send output
}

func renderDeviceGroups(w io.Writer, dg dsl.DeviceGroups) {
	table := NewTable(w)
	table.SetHeader([]string{"Device Group", "No of Devices"})
	for _, deviceGroup := range dg.DeviceGroup {
		l := "0"
		if deviceGroup.Devices != nil {
			l = strconv.Itoa(len(*deviceGroup.Devices))
		}
		table.Append([]string{deviceGroup.DeviceGroupName, l})
	}
	table.Render() // Send output
}

func summarizeInstallation(ci dsl.ConnectionInfo) (err error) {
	err = collectSystemDetails(resty.DefaultClient, ci, os.Stdout)
	if err != nil {
		return
	}

	df, err := collectDeviceFacts(resty.DefaultClient, ci, os.Stdout)
	if err != nil {
		return
	}

	renderDeviceTable(os.Stdout, df)

	dg, err := collectDeviceGroups(resty.DefaultClient, ci, os.Stdout)
	if err != nil {
		return
	}

	renderDeviceGroups(os.Stdout, dg)

	fmt.Println("")

	return
}

func init() {
	summarizeCmd.AddCommand(summarizeInstallationCmd)
}
