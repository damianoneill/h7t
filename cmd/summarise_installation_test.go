package cmd

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/damianoneill/h7t/dsl"
	"github.com/jarcoal/httpmock"
	"github.com/tj/assert"
	"gopkg.in/resty.v1"
)

var client = resty.New()

func Test_collectSystemDetails(t *testing.T) {
	type args struct {
		ci dsl.ConnectionInfo
	}
	tests := []struct {
		name               string
		args               args
		wantStdoutContains string
		wantErr            bool
	}{
		{
			name: "Test invalid system details",
			args: args{
				ci: dsl.ConnectionInfo{
					Authority: "localhost:8080",
					Username:  "root",
					Password:  "changeme",
				},
			},
			wantStdoutContains: "",
			wantErr:            true,
		},
		{
			name: "Test valid system details",
			args: args{
				ci: dsl.ConnectionInfo{
					Authority: "localhost:8080",
					Username:  "root",
					Password:  "changeme",
				},
			},
			wantStdoutContains: "2019-11-12T11:23:53Z",
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				// create valid resty mock
				httpmock.ActivateNonDefault(client.GetClient())
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", "https://localhost:8080/api/v1/system-details/",
					httpmock.NewStringResponder(200, `{"server-time": "2019-11-12T11:23:53Z","version": "HealthBot 2.1.0-beta"}`))

			}
			stdout := &bytes.Buffer{}
			if err := collectSystemDetails(client, tt.args.ci, stdout); (err != nil) != tt.wantErr {
				t.Errorf("collectSystemDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !strings.Contains(stdout.String(), tt.wantStdoutContains) {
					t.Errorf("collectSystemDetails() = %v, should contain %v", stdout.String(), tt.wantStdoutContains)
				}
			}
		})
	}
}

func Test_collectDeviceFacts(t *testing.T) {
	type args struct {
		ci dsl.ConnectionInfo
	}
	tests := []struct {
		name               string
		args               args
		wantStdoutContains string
		wantErr            bool
	}{
		{
			name: "Test invalid device facts",
			args: args{
				ci: dsl.ConnectionInfo{
					Authority: "localhost:8080",
					Username:  "root",
					Password:  "changeme",
				},
			},
			wantStdoutContains: "",
			wantErr:            true,
		},
		{
			name: "Test valid device facts",
			args: args{
				ci: dsl.ConnectionInfo{
					Authority: "localhost:8080",
					Username:  "root",
					Password:  "changeme",
				},
			},
			wantStdoutContains: "Managed Devices: 1",
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				// create valid resty mock
				httpmock.ActivateNonDefault(client.GetClient())
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder("GET", "https://localhost:8080/api/v1/devices/facts/",
					httpmock.NewStringResponder(200, `[{"device-id":"mx960-3","facts":{"fpc":[{"description":"MPC7E 3D MRATE-12xQSFPP-XGE-XLGE-CGE","model-number":"MPC7E-MRATE","name":"FPC 11","part-number":"750-056519","serial-number":"CAFR4421","version":"REV 36"}],"hostname":"mx960-3","junos-info":[{"last-reboot-reason":"0x4000:VJUNOS reboot","mastership-state":"master","model":"RE-S-2X00x6","name":"re0","status":"OK","up-time":"30 days, 5 hours, 59 minutes, 51 seconds","version-info":{"build":8,"major":[19,3],"minor":["1"],"type":"R"}}],"platform":"MX960","platform-info":[{"name":"re0","platform":"MX960"}],"product":"MX","release":"19.3R1.8","serial-number":"JN1233EF1AFA"}}]`))

			}
			stdout := &bytes.Buffer{}
			var df dsl.DeviceFacts
			var err error
			if df, err = collectDeviceFacts(client, tt.args.ci, stdout); (err != nil) != tt.wantErr {
				t.Errorf("collectDeviceFacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !strings.Contains(stdout.String(), tt.wantStdoutContains) {
					t.Errorf("collectDeviceFacts() = %v, should contain %v", stdout.String(), tt.wantStdoutContains)
				}
				assert.Len(t, df, 1, "should contain 1 device")
			}
		})
	}
}

// HelperLoadBytes allows you to use relative path testdata directory as a place
// to load and store your data
func HelperLoadBytes(tb testing.TB, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)     // nolint : gosec
	if err != nil {
		tb.Fatal(err)
	}
	return bytes
}

func Test_renderDeviceTable(t *testing.T) {

	b := HelperLoadBytes(t, "device-facts.json")
	df := dsl.DeviceFacts{}
	_ = df.Unmarshal(b)

	type args struct {
		df dsl.DeviceFacts
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{
			name:  "Should contain MX960",
			args:  args{df: df},
			wantW: "MX960     19.3R1.8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			renderDeviceTable(w, tt.args.df)
			if gotW := w.String(); !strings.Contains(gotW, tt.wantW) {
				t.Errorf("renderDeviceTable() = %v, should contain %v", gotW, tt.wantW)
			}
		})
	}
}
