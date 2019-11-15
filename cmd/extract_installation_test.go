package cmd

import "testing"

import "github.com/stretchr/testify/assert"

func Test_appendResourceBody(t *testing.T) {
	// fugly test
	backup := []byte("{\"device-group\":[{\"device-group-name\":\"ptp\",\"devices\":[\"mmx960-1\",\"mmx960-3\"]}],")
	responseBody := []byte("{\"device\":[{\"device-id\":\"mmx960-1\",\"host\":\"172.30.177.102\"}]}")
	expected := []byte("{\"device-group\":[{\"device-group-name\":\"ptp\",\"devices\":[\"mmx960-1\",\"mmx960-3\"]}],\"device\":[{\"device-id\":\"mmx960-1\",\"host\":\"172.30.177.102\"}],")

	err := appendResourceBody(&backup, responseBody)
	assert.Nil(t, err, "Should not return an error %v", err)
	assert.Equal(t, expected, backup, "should start with { and end with , and separate objects with a comma")
}
