package component

import (
	"encoding/json"

	"lib.hemtjan.st/v2/class/device"
	"lib.hemtjan.st/v2/platform"
)

type Settable interface {
	GetID() string
	GetPlatform() platform.Type
	GetDeviceClass() device.Class
}

type BaseComponent interface {
	GetBaseReference() *Base
}

var _ Settable = Base{}

// Base are fields any component must have.
type Base struct {
	Name                 string         `json:"name"`
	ID                   string         `json:"uniq_id"`
	Platform             platform.Type  `json:"p"`
	DeviceClass          device.Class   `json:"dev_cla"`
	BaseTopic            string         `json:"~,omitempty"`
	CommandTopic         string         `json:"cmd_t,omitempty"`
	Encoding             string         `json:"e,omitempty"`
	QoS                  uint           `json:"qos,omitempty"`
	StateTopic           string         `json:"stat_t,omitempty"`
	Availability         []Availability `json:"avty,omitempty"`
	AvailabilityMode     string         `json:"avty_mode,omitempty"`
	AvailabilityTemplate string         `json:"avty_tpl,omitempty"`
	AvailabilityTopic    string         `json:"avty_t,omitempty"`
}

func (b Base) GetID() string {
	return b.ID
}

func (b Base) GetPlatform() platform.Type {
	return b.Platform
}

func (b Base) GetDeviceClass() device.Class {
	return b.DeviceClass
}

func (b *Base) GetBaseReference() *Base {
	return b
}

// Generic holds any component.
//
// Once you've determined platform and class, use [Generic.As] to get the
// typed component.
type Generic struct {
	Base
	Raw json.RawMessage `json:"-"`
}

func (g *Generic) As(out json.Unmarshaler) error {
	return json.Unmarshal(g.Raw, out)
}

func (g Generic) MarshalJSON() ([]byte, error) {
	return g.Raw, nil
}

func (g *Generic) UnmarshalJSON(data []byte) error {
	type internal Generic
	var i internal
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	i.Raw = data
	*g = Generic(i)

	return nil
}

// Availability represents the availability.
type Availability struct {
	PayloadAvailable    string `json:"pl_avail,omitempty"`
	PayloadNotAvailable string `json:"pl_not_avail,omitempty"`
	Topic               string `json:"t"`
	Template            string `json:"val_tpl,omitempty"`
}
