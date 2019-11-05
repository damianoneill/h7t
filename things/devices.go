package things

// Devices - collection of Device
type Devices struct {
	Device []Device `json:"device"`
}

// Authentication - Collection type for Auth options
type Authentication struct {
	Password struct {
		Password *string `json:"password"`
		Username *string `json:"username"`
	} `json:"password,omitempty" yaml:"password,omitempty"`
}

// IAgent - configure the NETCONF port
type IAgent struct {
	Port int `json:"port"`
}

// OpenConfig - configure the Open Config port
type OpenConfig struct {
	Port int `json:"port"`
}

// V2 - configure the SNMP community string
type V2 struct {
	Community string `json:"community"`
}

// Snmp - configure the SNMP port or Community String
type Snmp struct {
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
	V2   *V2 `json:"v2,omitempty" yaml:"v2,omitempty"`
}

// Juniper - option to define the Operating system
type Juniper struct {
	OperatingSystem string `json:"operating-system" yaml:"operating-system"`
}

// Cisco - option to define the Operating system
type Cisco struct {
	OperatingSystem string `json:"operating-system" yaml:"operating-system"`
}

// Vendor - Configure the Vendor information
type Vendor struct {
	Juniper *Juniper `json:"juniper,omitempty" yaml:"juniper,omitempty"`
	Cisco   *Cisco   `json:"cisco,omitempty" yaml:"cisco,omitempty"`
}

// Device - info needed to Register a Device in Healthbot
type Device struct {
	DeviceID       string          `json:"device-id" yaml:"device-id"`
	Host           string          `json:"host"`
	SystemID       string          `json:"system-id,omitempty" yaml:"system-id,omitempty"`
	Authentication *Authentication `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	IAgent         *IAgent         `json:"iAgent,omitempty" yaml:"iAgent,omitempty"`
	OpenConfig     *OpenConfig     `json:"open-config,omitempty" yaml:"open-config,omitempty"`
	Snmp           *Snmp           `json:"snmp,omitempty" yaml:"snmp,omitempty"`
	Vendor         *Vendor         `json:"vendor,omitempty" yaml:"vendor,omitempty"`
}

// Unmarshal - tries to Unmarshal yaml first, then json into the Devices struct
func (d *Devices) Unmarshal(data []byte) error {
	return unmarshal(data, d)
}

// Unmarshal - tries to Unmarshal yaml first, then json into the Device struct
func (d *Device) Unmarshal(data []byte) error {
	return unmarshal(data, d)
}
