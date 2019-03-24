package server

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/hemtjanst/bibliotek/device"
)

// Transport is the server's transport
type Transport interface {
	// DeviceState receives a device.Info on announce, leave or update
	DeviceState() chan *device.State
	device.Transport
}

// Manager holds all the devices and deals with updating device
// state as it receives it from the transport
type Manager struct {
	transport Transport
	devices   map[string]Device
	waitingOn map[string]chan struct{}
	sync.RWMutex
}

// New creates a new Manager
func New(t Transport) *Manager {
	return &Manager{
		transport: t,
		devices:   map[string]Device{},
		waitingOn: map[string]chan struct{}{},
	}
}

// Start starts the Manager's main loop, conusming messages and
// reacting to changes in the system.Start
func (m *Manager) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case d := <-m.transport.DeviceState():
			switch d.Action {
			case device.DeleteAction:
				m.RemoveDevice(d.Topic)
			case device.LeaveAction:
				m.UnreachableDevice(d.Topic)
			case device.UpdateAction:
				if !m.HasDevice(d.Topic) {
					m.AddDevice(d.Device)
				} else {
					m.UpdateDevice(d.Device)
				}
			}
		}
	}
}

// HasDevice checks if a device is registered on the topic
func (m *Manager) HasDevice(topic string) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.devices[topic]
	return ok
}

// AddDevice adds a new device to the manager's devices
func (m *Manager) AddDevice(d *device.Info) {
	dev, err := NewDevice(d, m.transport)
	if err != nil {
		log.Printf("Failed to create device: %v", err)
		return
	}
	log.Printf("Device Created: %+v", d)
	m.Lock()
	defer m.Unlock()
	m.devices[d.Topic] = dev
	if ch, ok := m.waitingOn[d.Topic]; ok {
		close(ch)
		delete(m.waitingOn, d.Topic)
	}
}

// UpdateDevice updates an existing device with the new info
func (m *Manager) UpdateDevice(d *device.Info) {
	m.Lock()
	defer m.Unlock()

	dev := m.devices[d.Topic]

	updates, err := dev.update(d)
	if err != nil {
		log.Printf("Cannot update device %s: %s", dev.Id(), err)
		return
	}
	for _, upd := range updates {
		log.Printf("[%s] %s changed \"%s\" -> \"%s\" (%+v)",
			dev.Id(),
			upd.Field,
			upd.Old,
			upd.New,
			upd.FeatureInfo,
		)
	}

}

// RemoveDevice removes a device based on the topic name
// Removing a device that does not exist is a no-op
func (m *Manager) RemoveDevice(topic string) {
	m.Lock()
	defer m.Unlock()
	delete(m.devices, topic)
}

// UnreachableDevice marks the device specified by the topic
// as unreachable. This does not delete the device.
// Marking a non-existant device as unreachable is a no-op
func (m *Manager) UnreachableDevice(topic string) {
	m.Device(topic).setReachability(false)
}

// Device returns a device associated with the topic
func (m *Manager) Device(topic string) Device {
	if m.HasDevice(topic) {
		m.RLock()
		defer m.RUnlock()
		return m.devices[topic]
	}
	err := errors.New("device not found")
	return &FakeDevice{Topic: topic, Err: err}
}

// Devices returns all devices known to the manager
func (m *Manager) Devices() []Device {
	m.RLock()
	defer m.RUnlock()
	var devs []Device
	for _, dev := range m.devices {
		devs = append(devs, dev)
	}
	return devs
}

// WaitForDevice waits for a device on the specified topic to appear
// This blocks until either the device shows up, or the context is cancelled
func (m *Manager) WaitForDevice(ctx context.Context, topic string) Device {
	if m.HasDevice(topic) {
		m.RLock()
		defer m.RUnlock()
		return m.devices[topic]
	}
	m.Lock()

	ch, ok := m.waitingOn[topic]
	if !ok {
		ch = make(chan struct{})
		m.waitingOn[topic] = ch
	}

	m.Unlock()

	select {
	case <-ch:
		return m.Device(topic)
	case <-ctx.Done():
		return &FakeDevice{Err: errors.New("context cancelled"), Topic: topic}
	}
}
