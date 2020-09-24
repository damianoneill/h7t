package dsl

// DeviceFacts - Provides Device Facts
type DeviceFacts []struct {
	DeviceID string `json:"device-id" yaml:"device-id"`
	Facts    struct {
		Hostname  string `json:"hostname"`
		JunosInfo []struct {
			LastRebootReason string `json:"last-reboot-reason"`
			MastershipState  string `json:"mastership-state"`
			Model            string `json:"model"`
			Name             string `json:"name"`
			Status           string `json:"status"`
			UpTime           string `json:"up-time"`
		} `json:"junos-info"`
		Platform     string `json:"platform"`
		PlatformInfo []struct {
			Name     string `json:"name"`
			Platform string `json:"platform"`
		} `json:"platform-info"`
		Product      string `json:"product"`
		Release      string `json:"release"`
		SerialNumber string `json:"serial-number" yaml:"serial-number"`
	} `json:"facts,omitempty"`
}

// Unmarshal - tries to Unmarshal yaml first, then json into the DeviceFacts struct
func (d *DeviceFacts) Unmarshal(data []byte) error {
	return unmarshal(data, d)
}

// Path - resource path for DeviceFacts
func (d *DeviceFacts) Path() string {
	return "/api/v1/devices/facts/"
}

// Count - no of components within a thing
func (d *DeviceFacts) Count() int {
	return 1
}

// InnerThings - returns inner things or empty slice
func (d *DeviceFacts) InnerThings() []Thing {
	return []Thing{}
}
