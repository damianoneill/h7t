package dsl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"
)

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
	Count() int           // no of components within a thing
	InnerThings() []Thing // if Things is a aggregation of Things
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
	t, err := GetToken(rc, ci)
	if err != nil {
		return err
	}
	resp, err := rc.R().
		SetAuthToken(t.AccessToken).
		SetBasicAuth(ci.Username, ci.Password).
		Get("https://" + ci.Authority + thing.Path())
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.String())
	}
	err = thing.Unmarshal(resp.Body())
	return
}

// PostThingToResource - Use the Things Path function to build a POST REST requests and marshal body
func PostThingToResource(rc *resty.Client, thing Thing, ci ConnectionInfo, shouldCommit bool) (err error) {
	t, err := GetToken(rc, ci)
	if err != nil {
		return err
	}
	resp, err := rc.R().
		SetAuthToken(t.AccessToken).
		SetBasicAuth(ci.Username, ci.Password).
		SetBody(thing).
		Post("https://" + ci.Authority + thing.Path())
	if err != nil {
		return
	}
	switch resp.StatusCode() {
	case http.StatusOK:
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
	t, err := GetToken(rc, ci)
	if err != nil {
		return err
	}
	resp, err := rc.R().
		SetAuthToken(t.AccessToken).
		SetBasicAuth(ci.Username, ci.Password).
		Delete("https://" + ci.Authority + thing.Path())
	if err != nil {
		return
	}
	switch resp.StatusCode() {
	case http.StatusNoContent:
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

type Token struct {
	AccessToken string `json:"accessToken"`
	// FirstLogin          bool   `json:"firstLogin"`
	// RefreshToken        string `json:"refreshToken"`
	// RefreshTokenExpires int    `json:"refreshTokenExpires"`
	// TokenExpires        int    `json:"tokenExpires"`
}

// GetToken added for HB 3.0
func GetToken(rc *resty.Client, ci ConnectionInfo) (t *Token, err error) {
	r, err := rc.R().
		SetBody(`{"userName":"`+ci.Username+`", "password":"`+ci.Password+`"}`).
		SetHeader("Content-Type", "application/json").
		Post("https://" + ci.Authority + "/api/v1/login")

	if err != nil {
		return nil, fmt.Errorf("problem with authenticating to healthbot: %w", err)
	}

	var token Token
	err = json.Unmarshal(r.Body(), &token)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall authentication json response: %w", err)
	}

	// set token for future requests against this server
	return &token, err
}
