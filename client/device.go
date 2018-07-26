package client

import (
	"encoding/json"
	"github.com/hemtjanst/bibliotek/device"
	"github.com/hemtjanst/bibliotek/feature"
)

// Device in client package is used by applications that are talking to the actual devices,
// For example an application that's controlling lights over z-wave
type Device interface {
	// Common contains the common methods for client and server
	device.Common

	// Feature returns the feature with client functions (OnSet(), Update()) available.
	// Fetching a feature that doesn't exist will _NOT_ return nil, but instead create
	// a *feature.Fake. To check if feature actually exists; call Feature("name").Exists()
	Feature(name string) Feature

	// Features returns a slice of all available features
	Features() []Feature
}

// NewClient will create a new client device from the device.Info.
// It spawns off a goroutine that checks for inbound discover-messages
// and returns meta-data to the announce-topic.
// The transport is responsible for closing the channel, at that point the
// goroutine will stop.
func NewDevice(info *device.Info, transport device.Transport) (Device, error) {
	if info.LastWillID == "" {
		info.LastWillID = transport.LastWillID()
	}
	d := &clientDev{
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
	go func() {
		ch := transport.Discover()
		for {
			_, open := <-ch
			if !open {
				return
			}
			meta, _ := json.Marshal(d.Info)
			transport.PublishMeta(d.Id(), meta)
		}
	}()
	return d, nil
}

type clientDev struct {
	device.Device
}

func (d *clientDev) Feature(name string) Feature {
	return d.Device.Feature(name)
}

func (d *clientDev) Features() (fts []Feature) {
	for _, ft := range d.Device.Features {
		fts = append(fts, ft)
	}
	return
}
