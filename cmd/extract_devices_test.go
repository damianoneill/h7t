package cmd

import (
	"testing"

	"github.com/tj/assert"

	"github.com/spf13/afero"

	"github.com/damianoneill/h7t/pkg/dsl"
)

func TestWriteDevicesToFile(t *testing.T) {
	AppFs = afero.NewMemMapFs()

	type args struct {
		thing     dsl.Thing
		namedfile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "No devices should produce an error",
			args: args{
				thing: &dsl.Devices{
					Device: nil,
				},
				namedfile: "doesnt matter",
			},
			wantErr: true,
		},
		{
			name: "Valid devices should produce an file",
			args: args{
				thing: &dsl.Devices{
					Device: []dsl.Device{{
						DeviceID: "mx1",
						Host:     "10.0.0.1",
					}},
				},
				namedfile: "devices.yml",
			},
			wantErr: false,
		},
		{
			name: "Valid device-groups should produce an file",
			args: args{
				thing: &dsl.DeviceGroups{
					DeviceGroup: []dsl.DeviceGroup{{
						DeviceGroupName: "dg1",
					}},
				},
				namedfile: "./device-groups.yml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteThingsToFile(tt.args.thing, tt.args.namedfile); (err != nil) != tt.wantErr {
				t.Errorf("WriteThingsToFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				exists, existErr := afero.Exists(AppFs, tt.args.namedfile)
				assert.Nil(t, existErr, "Call to exist should not return an error")
				assert.True(t, exists, "File should exist")
			}
		})
	}
}
