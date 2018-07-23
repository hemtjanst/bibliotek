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
	info        *Info
	features    map[string]feature.Feature
	transporter DeviceTransporter
}

func (d *device) Id() string           { return d.info.Topic }
func (d *device) Name() string         { return d.info.Name }
func (d *device) Manufacturer() string { return d.info.Manufacturer }
func (d *device) Model() string        { return d.info.Model }
func (d *device) SerialNumber() string { return d.info.SerialNumber }
func (d *device) Type() string         { return d.info.Type }
func (d *device) Exists() bool         { return true }

func (d *device) Features() []feature.Feature {
	var fts []feature.Feature
	for _, ft := range d.features {
		fts = append(fts, ft)
	}
	return fts
}

func (d *device) Feature(name string) feature.Feature {
	if ft, ok := d.features[name]; ok {
		return ft
	}
	err := errors.New("feature not found")
	return &feature.Fake{FeatureName: name, Err: err}
}

type Device interface {
	Id() string
	Name() string
	Manufacturer() string
	Model() string
	SerialNumber() string
	Type() string
	Feature(name string) feature.Feature
	Exists() bool
	Features() []feature.Feature
}

type DeviceTransporter interface {
	Publish(topic string, payload []byte, retain bool)
	Subscribe(topic string) chan []byte
}

func New(info *Info, transporter DeviceTransporter) (Device, error) {
	if info == nil {
		return nil, fmt.Errorf("cannot create device without info")
	}
	if info.Topic == "" {
		return nil, fmt.Errorf("cannot have a device info with an empty topic")
	}

	d := &device{
		info:        info,
		features:    map[string]feature.Feature{},
		transporter: transporter,
	}

	if info.Features != nil {
		for name, ft := range info.Features {
			if ft.GetTopic == "" {
				ft.GetTopic = info.Topic + "/" + name + "/get"
			}
			if ft.SetTopic == "" {
				ft.SetTopic = info.Topic + "/" + name + "/set"
			}
			d.features[name] = feature.New(name, ft, d)
		}
	}

	return d, nil
}
