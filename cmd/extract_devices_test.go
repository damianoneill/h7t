package cmd

import (
	"testing"

	"github.com/tj/assert"

	"github.com/spf13/afero"

	"github.com/damianoneill/h7t/dsl"
)

func TestWriteDevicesToFile(t *testing.T) {

	AppFs = afero.NewMemMapFs()

	type args struct {
		devices   dsl.Devices
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
				devices: dsl.Devices{
					Device: nil,
				},
				namedfile: "doesnt matter",
			},
			wantErr: true,
		},
		{
			name: "Valid devices should produce an file",
			args: args{
				devices: dsl.Devices{
					Device: []dsl.Device{dsl.Device{
						DeviceID: "mx1",
						Host:     "10.0.0.1",
					}},
				},
				namedfile: "devices.yml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteDevicesToFile(tt.args.devices, tt.args.namedfile); (err != nil) != tt.wantErr {
				t.Errorf("WriteDevicesToFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				exists, err := afero.Exists(AppFs, "devices.yml")
				assert.Nil(t, err, "Call to exist should not return an error")
				assert.True(t, exists, "File should exist")
			}

		})
	}
}
