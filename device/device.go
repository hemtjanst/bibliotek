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

type device struct {
	info      *Info
	features  map[string]feature.Feature
	transport Transport
}

func (d *device) Id() string           { return d.info.Topic }
func (d *device) Name() string         { return d.info.Name }
func (d *device) Manufacturer() string { return d.info.Manufacturer }
func (d *device) Model() string        { return d.info.Model }
func (d *device) SerialNumber() string { return d.info.SerialNumber }
func (d *device) Type() string         { return d.info.Type }
func (d *device) Exists() bool         { return true }

func (d *device) getFeature(name string) feature.Feature {
	if ft, ok := d.features[name]; ok {
		return ft
	}
	err := errors.New("feature not found")
	return &feature.Fake{FeatureName: name, Err: err}
}

func create(d *device) error {
	if d.info == nil {
		return fmt.Errorf("cannot create device without info")
	}
	if d.info.Topic == "" {
		return fmt.Errorf("cannot have a device info with an empty topic")
	}

	if d.info.Features != nil {
		for name, ft := range d.info.Features {
			if ft.GetTopic == "" {
				ft.GetTopic = d.info.Topic + "/" + name + "/get"
			}
			if ft.SetTopic == "" {
				ft.SetTopic = d.info.Topic + "/" + name + "/set"
			}
			d.features[name] = feature.New(name, ft, d)
		}
	}

	return nil
}
