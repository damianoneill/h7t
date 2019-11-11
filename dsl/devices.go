package dsl

// Devices - collection of Device
type Devices struct {
	Device []Device `json:"device"`
}

// Password - wrapper for uname/password
type Password struct {
	Username *string `json:"username" csv:"username,omitempty"`
	Password *string `json:"password" csv:"password,omitempty"`
}

// Authentication - Collection type for Auth options
type Authentication struct {
	Password `json:"password,omitempty" yaml:"password,omitempty"`
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
	V2   *V2 `json:"v2,omitempty" yaml:"v2,omitempty"`
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
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
	DeviceID        string  `json:"device-id" yaml:"device-id" csv:"device-id"`
	Host            string  `json:"host" csv:"host"`
	SystemID        *string `json:"system-id,omitempty" yaml:"system-id,omitempty" csv:"-"`
	*Authentication `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	IAgent          *IAgent     `json:"iAgent,omitempty" yaml:"iAgent,omitempty" csv:"-"`
	OpenConfig      *OpenConfig `json:"open-config,omitempty" yaml:"open-config,omitempty" csv:"-"`
	Snmp            *Snmp       `json:"snmp,omitempty" yaml:"snmp,omitempty" csv:"-"`
	Vendor          *Vendor     `json:"vendor,omitempty" yaml:"vendor,omitempty" csv:"-"`
}

// Unmarshal - tries to Unmarshal yaml first, then json into the Devices struct
func (d *Devices) Unmarshal(data []byte) error {
	return unmarshal(data, d)
}

// Path - resource path for Devices
func (d *Devices) Path() string {
	return "/api/v1/devices/"
}

// Unmarshal - tries to Unmarshal yaml first, then json into the Device struct
func (d *Device) Unmarshal(data []byte) error {
	return unmarshal(data, d)
}

// Path - resource path for Device
func (d *Device) Path() string {
	return "/api/v1/device/" + d.DeviceID + "/"
}
