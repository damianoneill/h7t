package dsl

import (
	"encoding/json"
	"errors"
	"io"

	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"
)

// SystemDetails - Provides some basic hb info
type SystemDetails struct {
	ServerTime string `json:"server-time" yaml:"server-time"`
	Version    string `json:"version"`
}

// Unmarshal - tries to Unmarshal yaml first, then json into the Devices struct
func (d *SystemDetails) Unmarshal(data []byte) error {
	return unmarshal(data, d)
}

// Path - resource path for Devices
func (d *SystemDetails) Path() string {
	return "/api/v1/system-details/"
}

// ConnectionInfo - used to describe the resource endpoints
type ConnectionInfo struct {
	Authority string
	Username  string
	Password  string
}

// Thing - nouns that h7t act on
type Thing interface {
	Unmarshal(data []byte) error
	Path() string
}

func unmarshal(data []byte, t interface{}) error {
	if err := yaml.Unmarshal(data, t); err != nil {
		if err := json.Unmarshal(data, t); err != nil {
			return err
		}
	}
	return nil
}

// ReadThingFromFile - loads the file contents (yaml or json) into a thing, io.ReadFile passed so it can be UTed
func ReadThingFromFile(thing Thing, filename string, readfile func(filename string) ([]byte, error)) (err error) {
	b, err := readfile(filename)
	if err != nil {
		return
	}
	err = thing.Unmarshal(b)
	return
}

// WriteThingToFile - marshsal a thing to yaml and write to file
func WriteThingToFile(thing Thing, fw io.Writer) (err error) {
	b, err := yaml.Marshal(thing)
	if err != nil {
		return
	}
	_, err = fw.Write(b)
	return
}

// ExtractThingFromResource - Use the Things Path function to build a GET REST requests and unmarshal body to yaml
func ExtractThingFromResource(rc *resty.Client, thing Thing, ci ConnectionInfo) (err error) {
	resp, err := rc.R().
		SetBasicAuth(ci.Username, ci.Password).
		Get("https://" + ci.Authority + thing.Path())
	if err != nil {
		return
	}
	err = thing.Unmarshal(resp.Body())
	return
}

// PostThingToResource - Use the Things Path function to build a POST REST requests and marshal body
func PostThingToResource(rc *resty.Client, thing Thing, ci ConnectionInfo, shouldCommit bool) (err error) {
	resp, err := rc.R().
		SetBasicAuth(ci.Username, ci.Password).
		SetBody(thing).
		Post("https://" + ci.Authority + thing.Path())
	if err != nil {
		return
	}
	switch resp.StatusCode() {
	case 200:
		break
	default:
		return errors.New("Problem updating Thing: %v " + resp.String())
	}
	if shouldCommit {
		_, err = rc.R().
			SetBasicAuth(ci.Username, ci.Password).
			SetBody(thing).
			Post("https://" + ci.Authority + "/api/v1/configuration/")
		if err != nil {
			return
		}
	}
	return
}

// DeleteThingToResource - Use the Things Path function to build a DELETE REST request
func DeleteThingToResource(rc *resty.Client, thing Thing, ci ConnectionInfo, shouldCommit bool) (err error) {
	resp, err := rc.R().
		SetBasicAuth(ci.Username, ci.Password).
		Delete("https://" + ci.Authority + thing.Path())
	if err != nil {
		return
	}
	switch resp.StatusCode() {
	case 204:
		break
	default:
		return errors.New("Problem deleting Thing: %v " + resp.String())
	}
	if shouldCommit {
		_, err = rc.R().
			SetBasicAuth(ci.Username, ci.Password).
			SetBody(thing).
			Post("https://" + ci.Authority + "/api/v1/configuration/")
		if err != nil {
			return
		}
	}
	return
}
