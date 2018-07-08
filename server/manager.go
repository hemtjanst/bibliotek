package server

import (
	"context"
	"github.com/hemtjanst/bibliotek/device"
	"log"
)

type Transport interface {
	// HandleDeviceState should get a callback when a devices joins, leaves or
	// updates its metadata
	DeviceState() chan *device.DeviceInfo

	/*
		// HandleFeatureState is called with the id of the device and a feature name,
		// whenever that feature updates - the callback is called with the new value
		FeatureState(device, feature string) chan string

		// SetState should update the value of a device feature
		SetState(id, feature, newState string) error
	*/
}

type Manager struct {
	transport Transport
}

func New(t Transport) *Manager {
	return &Manager{t}
}

func (m *Manager) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case d := <-m.transport.DeviceState():
			log.Printf("Device State: %v", d)
		}
	}
}
