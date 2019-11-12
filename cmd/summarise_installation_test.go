package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/damianoneill/h7t/dsl"
	"github.com/jarcoal/httpmock"
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
				if !strings.Contains(stdout.String(), "2019-11-12T11:23:53Z") {
					t.Errorf("collectSystemDetails() = %v, should contain %v", stdout.String(), tt.wantStdoutContains)
				}
			}
		})
	}
}
