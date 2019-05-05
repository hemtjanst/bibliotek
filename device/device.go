package device

import (
	"errors"
	"fmt"
	"strconv"

	"lib.hemtjan.st/feature"
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

type InfoUpdate struct {
	Field       string
	Old         string
	New         string
	FeatureInfo []*feature.InfoUpdate
}

func (d *Device) UpdateInfo(info *Info) (updates []*InfoUpdate, err error) {

	if d.Info.Topic != info.Topic {
		err = errors.New("attempted to update device topic")
		return
	}

	if d.Info.LastWillID != info.LastWillID {
		updates = append(updates, &InfoUpdate{Field: "lastWillID", Old: d.Info.LastWillID, New: info.LastWillID})
		d.Info.LastWillID = info.LastWillID
	}

	if d.Info.Name != info.Name {
		updates = append(updates, &InfoUpdate{Field: "name", Old: d.Info.Name, New: info.Name})
		d.Info.Name = info.Name
	}

	if d.Info.Reachable != info.Reachable {
		nv := "0"
		ov := "0"
		if d.Info.Reachable {
			ov = "1"
		}
		if info.Reachable {
			nv = "1"
		}
		updates = append(updates, &InfoUpdate{Field: "reachable", Old: ov, New: nv})
		d.Info.Reachable = info.Reachable
	}

	if d.Info.Type != info.Type {
		updates = append(updates, &InfoUpdate{Field: "type", Old: d.Info.Type, New: info.Type})
		d.Info.Type = info.Type
	}

	if d.Info.Manufacturer != info.Manufacturer {
		updates = append(updates, &InfoUpdate{Field: "manufacturer", Old: d.Info.Manufacturer, New: info.Manufacturer})
		d.Info.Manufacturer = info.Manufacturer
	}

	if d.Info.Model != info.Model {
		updates = append(updates, &InfoUpdate{Field: "model", Old: d.Info.Model, New: info.Model})
		d.Info.Model = info.Model
	}

	if d.Info.SerialNumber != info.SerialNumber {
		updates = append(updates, &InfoUpdate{Field: "serialNumber", Old: d.Info.SerialNumber, New: info.SerialNumber})
		d.Info.SerialNumber = info.SerialNumber
	}

	for name, ft := range info.Features {
		if ft.GetTopic == "" {
			ft.GetTopic = d.Info.Topic + "/" + name + "/get"
		}
		if ft.SetTopic == "" {
			ft.SetTopic = d.Info.Topic + "/" + name + "/set"
		}
		curFt := d.Feature(name)
		if !curFt.Exists() {
			d.Features[name] = feature.New(name, ft, d)
			updates = append(updates, &InfoUpdate{
				Field: "feature",
				Old:   "",
				New:   name,
				FeatureInfo: []*feature.InfoUpdate{
					{Name: "min", New: strconv.Itoa(ft.Min)},
					{Name: "max", New: strconv.Itoa(ft.Max)},
					{Name: "step", New: strconv.Itoa(ft.Step)},
					{Name: "getTopic", New: ft.GetTopic},
					{Name: "setTopic", New: ft.SetTopic},
				},
			})
		} else {
			ftUpd := curFt.UpdateInfo(ft)
			if len(ftUpd) > 0 {
				updates = append(updates, &InfoUpdate{
					Field:       "feature",
					Old:         name,
					New:         name,
					FeatureInfo: ftUpd,
				})
				for _, uu := range ftUpd {
					if uu.Name == "getTopic" || uu.Name == "setTopic" {
						d.Transport.Resubscribe(uu.Old, uu.New)
					}
				}
			}

		}
	}
	for name, ft := range d.Features {
		if _, ok := info.Features[name]; !ok {
			updates = append(updates, &InfoUpdate{
				Field: "feature",
				Old:   name,
				New:   "",
			})
			d.Transport.Unsubscribe(ft.GetTopic())
			d.Transport.Unsubscribe(ft.SetTopic())
			delete(d.Features, name)
		}
	}

	return
}
