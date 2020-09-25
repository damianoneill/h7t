package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/damianoneill/h7t/pkg/dsl"
	"github.com/damianoneill/net/netconf"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var verbose string

func sendRPC(d *dsl.Device, b []byte) (err error) {
	sshConfig := &ssh.ClientConfig{
		User:            *d.Password.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(*d.Password.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint:gosec
	}

	port := 830
	if d.IAgent != nil {
		port = d.IAgent.Port
	}

	serverAddress := fmt.Sprintf("%v:%d", d.Host, port)
	s, err := netconf.NewRPCSession(context.Background(), sshConfig, serverAddress)
	if err != nil {
		return
	}
	reply, err := s.Execute(netconf.Request(string(b)))
	s.Close()

	if verbose == "true" {
		fmt.Fprintf(os.Stdout, "rpc response: %v", reply.Data)
	}

	return
}

func applyNetconfRPC(devicesFiles []string, rpcFilename string) (err error) {
	rpcFile, err := afero.ReadFile(AppFs, rpcFilename)
	if err != nil {
		return
	}
	for _, device := range devicesFiles {
		f, deviceErr := afero.ReadFile(AppFs, device)
		if deviceErr != nil {
			return
		}
		devices := dsl.Devices{}
		err = devices.Unmarshal(f)
		if err != nil {
			return
		}
		for _, d := range devices.Device {
			rpcError := sendRPC(&d, rpcFile) // nolint:scopelint
			if rpcError != nil {
				fmt.Fprintf(os.Stdout, "Problem with Device %v: %v", d.DeviceID, rpcError)
				// do not error out, could be device is gone, continue trying others
			}
		}
		fmt.Fprintf(os.Stdout, "Sent request to %v Devices defined in %v\n", devices.Count(), device)
	}
	return
}

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Load Devices with configuration",
	Long: `Load into Devices, configuration defined in the netconf rpc file.

E.g. 

$ cat sample.rpc 
<edit-config>
  <target>
    <candidate />
  </target>
  <config>
    <configuration>
      <system>
        <login>
          <message>Login Message...</message>
        </login>
      </system>
    </configuration>
  </config>
</edit-config>
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		verbose = cmd.Flag("verbose").Value.String()
		devices, err := getDirectoryContents(cmd.Flag("input_directory").Value.String())
		if err != nil {
			return
		}
		err = applyNetconfRPC(devices, cmd.Flag("netconf_rpc").Value.String())
		return
	},
}

func init() {
	configureCmd.AddCommand(devicesCmd)
	devicesCmd.PersistentFlags().StringP("input_directory", "i", ".", "directory where the device configuration will be loaded from")
	devicesCmd.PersistentFlags().StringP("netconf_rpc", "f", "", "file that contains a netconf rpc")
}
