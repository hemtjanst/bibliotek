package server

import (
	"github.com/hemtjanst/bibliotek/device"
	"github.com/hemtjanst/bibliotek/feature"
)

// Device is used by applications that monitor or send commands to the
// different client devices.
// For examples an application that integrates with a personal assistant
type Device interface {
	// Common contains the common methods for client and server
	device.Common

	// Feature returns the feature with server functions (Set(), OnUpdate()) available.
	// Fetching a feature that doesn't exist will _NOT_ return nil, but instead create
	// a *feature.Fake. To check if feature actually exists; call Feature("name").Exists()
	Feature(name string) Feature

	// Exists returns true if the device exists
	Exists() bool

	// Features returns a slice of available features
	Features() []Feature

	// Changes the reachable attribute of a device
	setReachability(bool)

	// Updates the device with the new representation
	update(*device.Info) ([]*device.InfoUpdate, error)
}

// NewDevice should normally only be called with data from announcements.
func NewDevice(info *device.Info, transport device.Transport) (Device, error) {
	d := &serverDev{
		Device: device.Device{
			Info:      info,
			Features:  map[string]feature.Feature{},
			Transport: transport,
		},
	}
	err := device.Create(&d.Device)
	if err != nil {
		return nil, err
	}
	return d, nil
}

type serverDev struct {
	device.Device
}

func (d *serverDev) Feature(name string) Feature {
	return d.Device.Feature(name)
}

func (d *serverDev) Features() (fts []Feature) {
	for _, ft := range d.Device.Features {
		fts = append(fts, ft)
	}
	return
}

func (d *serverDev) setReachability(r bool) {
	d.Info.Reachable = r
}

// update updates the device representation
func (d *serverDev) update(info *device.Info) ([]*device.InfoUpdate, error) {
	return d.Device.UpdateInfo(info)
}
