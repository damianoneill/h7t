package things

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

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
	var device Device
	err := device.Unmarshal([]byte{45})
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
