package device

import (
	"errors"
	"fmt"

	"github.com/hemtjanst/bibliotek/feature"
)

type Info struct {
	Topic        string                   `json:"topic"`
	Name         string                   `json:"name"`
	Manufacturer string                   `json:"manufacturer"`
	Model        string                   `json:"model"`
	SerialNumber string                   `json:"serialNumber"`
	Type         string                   `json:"type"`
	LastWillID   string                   `json:"lastWillID,omitempty"`
	Features     map[string]*feature.Info `json:"feature"`
	Reachable    bool                     `json:"-"`
}

type Device struct {
	Info      *Info
	Features  map[string]feature.Feature
	Transport Transport
}

func (d *Device) Id() string           { return d.Info.Topic }
func (d *Device) Name() string         { return d.Info.Name }
func (d *Device) Manufacturer() string { return d.Info.Manufacturer }
func (d *Device) Model() string        { return d.Info.Model }
func (d *Device) SerialNumber() string { return d.Info.SerialNumber }
func (d *Device) Type() string         { return d.Info.Type }
func (d *Device) Exists() bool         { return true }

func (d *Device) Feature(name string) feature.Feature {
	if ft, ok := d.Features[name]; ok {
		return ft
	}
	err := errors.New("feature not found")
	return &feature.Fake{FeatureName: name, Err: err}
}

func Create(d *Device) error {
	if d.Info == nil {
		return fmt.Errorf("cannot create device without info")
	}
	if d.Info.Topic == "" {
		return fmt.Errorf("cannot have a device info with an empty topic")
	}

	if d.Info.Features != nil {
		for name, ft := range d.Info.Features {
			if ft.GetTopic == "" {
				ft.GetTopic = d.Info.Topic + "/" + name + "/get"
			}
			if ft.SetTopic == "" {
				ft.SetTopic = d.Info.Topic + "/" + name + "/set"
			}
			d.Features[name] = feature.New(name, ft, d)
		}
	}

	return nil
}
