package device

import (
	"lib.hemtjan.st/v2/component"
)

// Device represents a device.
//
// A device groups multilpe entities together.
type Device struct {
	Info       Info                          `json:"dev"`
	Origin     Origin                        `json:"o,omitempty"`
	Components map[string]component.Settable `json:"cmps,omitempty"`
}

func (d *Device) DiscoveryTopic() string {
	return "homeassistant/device/" + d.Info.ID + "/config"
}

func (d *Device) SetComponent(name string, comp component.Settable) error {
	if d.Components == nil {
		d.Components = map[string]component.Settable{}
	}

	d.Components[name] = comp
	return nil
}
