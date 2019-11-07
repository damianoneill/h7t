package dsl

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"
)

var ci = ConnectionInfo{
	Authority: "localhost:8080",
	Username:  "root",
	Password:  "Be1fast",
}

var fakeDevices = Devices{
	Device: []Device{{
		DeviceID: "test-device",
		Host:     "10.0.0.1",
	}},
}

var client = resty.New()

type FakeReadFiler struct {
	Str string
}

// here's a fake ReadFile method that matches the signature of ioutil.ReadFile
func (f FakeReadFiler) ReadFile(filename string) ([]byte, error) {
	buf := bytes.NewBufferString(f.Str)
	return ioutil.ReadAll(buf)
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

func TestUnmarshalFailure(t *testing.T) {
	var devices Devices
	err := devices.Unmarshal([]byte{45})
	assert.EqualError(t, err, "invalid character ' ' in numeric literal", "should have generated an error")
}

func TestDeviceYamlParsing(t *testing.T) {
	var devices Devices
	err := devices.Unmarshal(HelperLoadBytes(t, "./devices.yml"))
	assert.Nil(t, err, "Failed to parse yaml representation of Devices")
	assert.Len(t, devices.Device, 3, "Expected to parse 3 Devices")
	assert.EqualValues(t, "4200_1", devices.Device[0].DeviceID, "Yaml type with a hyphen, didn't decode correctly")
}

func TestDeviceOmit(t *testing.T) {
	var devices Devices
	_ = devices.Unmarshal(HelperLoadBytes(t, "./devices.yml"))
	maxDevice, err := yaml.Marshal(devices.Device[0])
	if err != nil {
		assert.Nil(t, err, "Failed to marshal devices to yaml")
	}
	assert.NotContains(t, string(maxDevice), "cisco", "cisco should be ignored")

	minDevice, err := json.Marshal(devices.Device[2])
	if err != nil {
		assert.Nil(t, err, "Failed to marshal devices to json")
	}
	assert.NotContains(t, string(minDevice), "authentication", "Optional Authentication type was not ignored")
	assert.NotContains(t, string(minDevice), "iAgent", "Optional NETCONF type was not ignored")
	assert.NotContains(t, string(minDevice), "open-config", "Optional Open Config type was not ignored")
	assert.NotContains(t, string(minDevice), "snmp", "Optional SNMP type was not ignored")
	assert.NotContains(t, string(minDevice), "vendor", "Optional vendor type was not ignored")
	partialDevice, err := json.Marshal(devices.Device[1])
	if err != nil {
		assert.Nil(t, err, "Failed to marshal devices to json")
	}
	assert.NotContains(t, string(partialDevice), "v2", "Optional Community type was not ignored")
	assert.NotContains(t, string(partialDevice), "juniper", "Optional Vendor juniper was not ignored")
}

func TestWriteThingToFile(t *testing.T) {
	devices := Devices{
		Device: []Device{{
			DeviceID: "localhost",
		}},
	}
	type args struct {
		thing Thing
	}
	tests := []struct {
		name     string
		args     args
		contains string
		wantErr  bool
	}{
		{"test writing a yaml thing", args{thing: &devices}, "device-id: localhost", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fw := &bytes.Buffer{}
			if err := WriteThingToFile(tt.args.thing, fw); (err != nil) != tt.wantErr {
				t.Errorf("WriteThingToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFw := fw.String(); !strings.Contains(gotFw, tt.contains) {
				t.Errorf("WriteThingToFile() = %v, want %v", gotFw, tt.contains)
			}
		})
	}
}

func TestExtractThingFromResource(t *testing.T) {

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://localhost:8080/api/v1/devices/",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, Devices{})
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	thing := Devices{}
	err := ExtractThingFromResource(client, &thing, ci)
	assert.Nil(t, err, "response from ExtractThingFromResource should be valid")

	assert.True(t, httpmock.GetTotalCallCount() == 1, "Expected Single call to Resource")

	_, isPresent := httpmock.GetCallCountInfo()["GET https://localhost:8080/api/v1/devices/"]
	assert.True(t, isPresent, "Should contain a correctly formatted GET")

}

func TestReadThingFromFile(t *testing.T) {
	b, _ := yaml.Marshal(fakeDevices)
	fake := FakeReadFiler{Str: string(b)}

	type args struct {
		thing    Thing
		filename string
		readfile func(filename string) ([]byte, error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Ensure valid Devices are returned",
			args{thing: &Devices{}, filename: "devices.yml", readfile: fake.ReadFile},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadThingFromFile(tt.args.thing, tt.args.filename, tt.args.readfile); (err != nil) != tt.wantErr {
				t.Errorf("ReadThingFromFile() error = %v, wantErr %v", err, tt.wantErr)
				assert.Len(t, tt.args.thing.(*Devices).Device, 1)
			}
		})
	}
}

func TestPostThingToResource(t *testing.T) {

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://localhost:8080/api/v1/devices/",
		httpmock.NewStringResponder(200, ``))

	type args struct {
		rc           *resty.Client
		thing        Thing
		ci           ConnectionInfo
		shouldCommit bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid post thing to resource", args{rc: client, thing: &fakeDevices, ci: ci}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostThingToResource(tt.args.rc, tt.args.thing, tt.args.ci, tt.args.shouldCommit); (err != nil) != tt.wantErr {
				t.Errorf("PostThingToResource() error = %v, wantErr %v", err, tt.wantErr)
				assert.Len(t, httpmock.GetTotalCallCount(), 1, "Expected Single call to Resource")

				_, isPresent := httpmock.GetCallCountInfo()["POST https://localhost:8080/api/v1/devices/"]
				assert.True(t, isPresent, "Should contain a correctly formatted POST")
			}
		})
	}
}

func TestDeleteThingToResource(t *testing.T) {

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://localhost:8080/api/v1/device/test-device/",
		httpmock.NewStringResponder(204, ``))

	type args struct {
		rc           *resty.Client
		thing        Thing
		ci           ConnectionInfo
		shouldCommit bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid delete thing to resource", args{rc: client, thing: &fakeDevices.Device[0], ci: ci}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteThingToResource(tt.args.rc, tt.args.thing, tt.args.ci, tt.args.shouldCommit); (err != nil) != tt.wantErr {
				t.Errorf("DeleteThingToResource() error = %v, wantErr %v", err, tt.wantErr)
				assert.Len(t, httpmock.GetTotalCallCount(), 1, "Expected Single call to Resource")

				_, isPresent := httpmock.GetCallCountInfo()["DELETE https://localhost:8080/api/v1/devices/test-device/"]
				assert.True(t, isPresent, "Should contain a correctly formatted POST")
			}
		})
	}
}
