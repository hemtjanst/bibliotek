package device

import "github.com/hemtjanst/bibliotek/feature"

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
	info     *Info
	features map[string]feature.Feature
}

func (d *device) Id() string           { return d.info.Topic }
func (d *device) Name() string         { return d.info.Name }
func (d *device) Manufacturer() string { return d.info.Manufacturer }
func (d *device) Model() string        { return d.info.Model }
func (d *device) SerialNumber() string { return d.info.SerialNumber }
func (d *device) Type() string         { return d.info.Type }

func (d *device) Features() map[string]feature.Feature {
	return d.features
}

func (d *device) Feature(name string) feature.Feature {
	if ft, ok := d.features[name]; ok {
		return ft
	}
	return nil
}

type Device interface {
	Id() string
	Name() string
	Manufacturer() string
	Model() string
	SerialNumber() string
	Type() string
	Feature(name string) feature.Feature
}

func New(info *Info) Device {
	if info == nil {
		info = &Info{}
	}
	d := &device{
		info:     info,
		features: map[string]feature.Feature{},
	}
	if info.Features != nil {
		for name, ft := range info.Features {
			d.features[name] = feature.New(name, ft)
		}
	}

	return d
}
