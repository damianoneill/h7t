package dsl

// SystemDetails - Provides some basic hb info
type SystemDetails struct {
	ServerTime string `json:"server-time" yaml:"server-time"`
	Version    string `json:"version"`
}

// Unmarshal - tries to Unmarshal yaml first, then json into the Devices struct
func (sd *SystemDetails) Unmarshal(data []byte) error {
	return unmarshal(data, sd)
}

// Path - resource path for Devices
func (sd *SystemDetails) Path() string {
	return "/api/v1/system-details/"
}

// Count - no of components within a thing
func (sd *SystemDetails) Count() int {
	return 1
}

// InnerThings - returns inner things or empty slice
func (sd *SystemDetails) InnerThings() []Thing {
	return []Thing{}
}
