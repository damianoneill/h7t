package dsl

// DeviceGroups - collection of Device Groups
type DeviceGroups struct {
	DeviceGroup []DeviceGroup `json:"device-group" yaml:"device-group"`
}

// DGAuthentication - Option to Override the individual Device Username/Passwords
type DGAuthentication struct {
	Password struct {
		Password *string `json:"password"`
		Username *string `json:"username"`
	} `json:"password,omitempty" yaml:"password,omitempty"`
}

// NativeGpb - Override the default JTI Port(s)
type NativeGpb struct {
	Ports []int `json:"ports"`
}

// DeviceGroup - info needed to Register a DeviceGroup in Healthbot
type DeviceGroup struct {
	DeviceGroupName string            `json:"device-group-name" yaml:"device-group-name"`
	Description     *string           `json:"description,omitempty" yaml:"description,omitempty"`
	Devices         *[]string         `json:"devices,omitempty" yaml:"devices,omitempty"`
	Playbooks       *[]string         `json:"playbooks,omitempty" yaml:"playbooks,omitempty"`
	Authentication  *DGAuthentication `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	NativeGpb       *NativeGpb        `json:"native-gpb,omitempty" yaml:"native-gpb,omitempty"`
}

// Unmarshal - tries to Unmarshal yaml first, then json into the DeviceGroups struct
func (d *DeviceGroups) Unmarshal(data []byte) error {
	return unmarshal(data, d)
}

// Path - resource path for DeviceGroups
func (d *DeviceGroups) Path() string {
	return "/api/v1/device-groups/"
}

// Count - no of components within a thing
func (d *DeviceGroups) Count() int {
	return len(d.DeviceGroup)
}

// InnerThings - returns inner things or empty slice
func (d *DeviceGroups) InnerThings() []Thing {
	b := make([]Thing, len(d.DeviceGroup))
	for i := range d.DeviceGroup {
		b[i] = &d.DeviceGroup[i]
	}
	return b
}

// Unmarshal - tries to Unmarshal yaml first, then json into the DeviceGroup struct
func (d *DeviceGroup) Unmarshal(data []byte) error {
	return unmarshal(data, d)
}

// Path - resource path for DeviceGroup
func (d *DeviceGroup) Path() string {
	return "/api/v1/device-group/" + d.DeviceGroupName + "/"
}

// Count - no of components within a thing
func (d *DeviceGroup) Count() int {
	return 1
}

// InnerThings - returns inner things or empty slice
func (d *DeviceGroup) InnerThings() []Thing {
	return []Thing{}
}
