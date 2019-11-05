package things

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Thing - nouns that h7t act on
type Thing interface {
	Unmarshal(data []byte) error
}

func unmarshal(data []byte, t interface{}) error {
	if err := yaml.Unmarshal(data, t); err != nil {
		if err := json.Unmarshal(data, t); err != nil {
			return err
		}
	}
	return nil
}
