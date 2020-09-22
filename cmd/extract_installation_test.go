package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_appendResourceBody(t *testing.T) {
	// fugly test
	backup := []byte("{\"device-group\":[{\"device-group-name\":\"ptp\",\"devices\":[\"mmx960-1\",\"mmx960-3\"]}],")
	responseBody := []byte("{\"device\":[{\"device-id\":\"mmx960-1\",\"host\":\"172.30.177.102\"}]}")
	expected := []byte("{\"device-group\":[{\"device-group-name\":\"ptp\",\"devices\":[\"mmx960-1\",\"mmx960-3\"]}],\"device\":[{\"device-id\":\"mmx960-1\",\"host\":\"172.30.177.102\"}],") // nolint:lll

	err := appendResourceBody(&backup, responseBody)
	assert.Nil(t, err, "Should not return an error %v", err)
	assert.Equal(t, expected, backup, "should start with { and end with , and separate objects with a comma")
}

func Test_removeBracketsFromJSONObject(t *testing.T) {
	b, err := removeBracketsFromJSONObject(HelperLoadBytes(t, "devices.json"))
	if err != nil {
		t.Errorf("removeBracketsFromJSONObject() error = %v", err)
		return
	}
	got := string(b)
	assert.True(t, strings.HasPrefix(got, "\"device\":"), "Should start with device colon")
	assert.True(t, strings.HasSuffix(got, "]"), "Should end with square bracket")
}
