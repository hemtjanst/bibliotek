package server

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/google/go-cmp/cmp"
	"github.com/hemtjanst/bibliotek/device"
)

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
			log.Printf("Device State: %+v", d)
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

func (m *Manager) HasDevice(topic string) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.devices[topic]
	return ok
}

func (m *Manager) AddDevice(d *device.Info) {
	dev, err := NewDevice(d, m.transport)
	if err != nil {
		log.Printf("Failed to create device: %v", err)
		return
	}
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
	upd, err := NewDevice(d, m.transport)
	if err != nil {
		log.Printf("Failed to create device: %v", err)
		return
	}
	m.Lock()
	defer m.Unlock()

	dev := m.devices[d.Topic]

	if upd.Name() != dev.Name() {
		log.Printf("Device has different Name: current %s, new %s", dev.Name(), upd.Name())
		return
	}
	if upd.Manufacturer() != dev.Manufacturer() {
		log.Printf("Device has different Manufacturer: current %s, new %s", dev.Manufacturer(), upd.Manufacturer())
		return
	}
	if upd.Model() != dev.Model() {
		log.Printf("Device has different Model: current %s, new %s", dev.Model(), upd.Model())
		return
	}
	if upd.SerialNumber() != dev.SerialNumber() {
		log.Printf("Device has different SerialNumber: current %s, new %s", dev.SerialNumber(), upd.SerialNumber())
		return
	}
	if upd.Type() != dev.Type() {
		log.Printf("Device has different Type: current %s, new %s", dev.Type(), upd.Type())
		return
	}

	oldft := dev.Features()
	newft := upd.Features()
	if diff := cmp.Diff(oldft, newft); diff != "" {
		log.Printf("Device has different features (-current +new):\n%s", diff)
	}

	dev.update(d)
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

func (m *Manager) Device(topic string) Device {
	if m.HasDevice(topic) {
		m.RLock()
		defer m.RUnlock()
		return m.devices[topic]
	}
	err := errors.New("device not found")
	return &FakeDevice{Topic: topic, Err: err}
}

func (m *Manager) Devices() []Device {
	m.RLock()
	defer m.RUnlock()
	var devs []Device
	for _, dev := range m.devices {
		devs = append(devs, dev)
	}
	return devs
}

func (m *Manager) WaitForDevice(topic string, ctx context.Context) Device {
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
