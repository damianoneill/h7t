package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/tj/assert"

	"github.com/damianoneill/h7t/pkg/dsl"
)

type FakeReadFiler struct {
	Str string
}

// here's a fake ReadFile method that matches the signature of ioutil.ReadFile
func (f FakeReadFiler) ReadFile(filename string) ([]byte, error) {
	buf := bytes.NewBufferString(f.Str)
	return ioutil.ReadAll(buf)
}

func TestReadCsvFile(t *testing.T) {
	fake := FakeReadFiler{Str: "device-id,host,username,password\nmx1,10.0.0.1,,"}

	type args struct {
		filename string
		readfile func(filename string) ([]byte, error)
	}
	tests := []struct {
		name        string
		args        args
		wantDevices []dsl.Device
		wantErr     bool
	}{
		{
			name:        "Test csv parsed correctly from file",
			args:        args{filename: "filename", readfile: fake.ReadFile},
			wantDevices: []dsl.Device{{DeviceID: "mx1", Host: "10.0.0.1"}},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDevices, err := ReadCsvFile(tt.args.filename, tt.args.readfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCsvFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, gotDevices[0].Host, tt.wantDevices[0].Host, "Host value should be the same")
			assert.Equal(t, gotDevices[0].DeviceID, tt.wantDevices[0].DeviceID, "Device ID value should be the same")
		})
	}
}
