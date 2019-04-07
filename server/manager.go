package server

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/hemtjanst/bibliotek/device"
)

// Transport is the server's transport
type Transport interface {
	// DeviceState receives a device.Info on announce, leave or update
	DeviceState() chan *device.State
	device.Transport
}

type Handler interface {
	AddedDevice(Device)
	UpdatedDevice(Device, []*device.InfoUpdate)
	RemovedDevice(Device)
}

type UpdateType string

const (
	AddedDevice   = "added"
	UpdatedDevice = "updated"
	RemovedDevice = "removed"
)

type Update struct {
	Type    UpdateType
	Device  Device
	Changes []*device.InfoUpdate
}

// Manager holds all the devices and deals with updating device
// state as it receives it from the transport
type Manager struct {
	updateChan chan<- Update
	transport  Transport
	devices    map[string]Device
	waitingOn  map[string]chan struct{}
	lock       sync.RWMutex
}

// New creates a new Manager
func New(t Transport) *Manager {
	return &Manager{
		transport: t,
		devices:   map[string]Device{},
		waitingOn: map[string]chan struct{}{},
	}
}

// SetUpdateChannel registers a channel that will receive Update when a device
// is added, removed or changed.
func (m *Manager) SetUpdateChannel(ch chan<- Update) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.updateChan != nil {
		close(m.updateChan)
	}
	m.updateChan = ch
}

// SetHandler is a wrapper around SetUpdateChannel for use with callbacks instead of channels,
// the handler is called synchronously
func (m *Manager) SetHandler(handler Handler) {
	ch := make(chan Update, 32)
	go func() {
		for {
			u, open := <-ch
			if !open {
				return
			}
			switch u.Type {
			case AddedDevice:
				handler.AddedDevice(u.Device)
			case UpdatedDevice:
				handler.UpdatedDevice(u.Device, u.Changes)
			case RemovedDevice:
				handler.RemovedDevice(u.Device)
			}
		}
	}()
	m.SetUpdateChannel(ch)
}

// Start starts the Manager's main loop, conusming messages and
// reacting to changes in the system.Start
func (m *Manager) Start(ctx context.Context) error {
	defer func() {
		m.lock.Lock()
		defer m.lock.Unlock()
		if m.updateChan != nil {
			close(m.updateChan)
			m.updateChan = nil
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case d, open := <-m.transport.DeviceState():
			if !open {
				return nil
			}
			switch d.Action {
			case device.DeleteAction:
				m.removeDevice(d.Topic)
			case device.LeaveAction:
				m.unreachableDevice(d.Topic)
			case device.UpdateAction:
				if !m.HasDevice(d.Topic) {
					m.addDevice(d.Device)
				} else {
					m.updateDevice(d.Device)
				}
			}
		}
	}
}

// HasDevice checks if a device is registered on the topic
func (m *Manager) HasDevice(topic string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.devices[topic]
	return ok
}

// AddDevice adds a new device to the manager's devices
func (m *Manager) addDevice(d *device.Info) {
	dev, err := NewDevice(d, m.transport)
	if err != nil {
		log.Printf("Failed to create device: %v", err)
		return
	}
	m.lock.Lock()
	m.devices[d.Topic] = dev

	var waitCh chan struct{}
	var ok bool

	if waitCh, ok = m.waitingOn[d.Topic]; ok {
		delete(m.waitingOn, d.Topic)
	}

	m.lock.Unlock()

	if waitCh != nil {
		close(waitCh)
	}

	if m.updateChan != nil {
		m.updateChan <- Update{Type: AddedDevice, Device: dev}
	}
}

// UpdateDevice updates an existing device with the new info
func (m *Manager) updateDevice(d *device.Info) {
	m.lock.Lock()
	dev := m.devices[d.Topic]
	changes, err := dev.update(d)
	m.lock.Unlock()
	if err != nil {
		log.Printf("Cannot update device %s: %s", dev.Id(), err)
		return
	}
	if len(changes) > 0 && m.updateChan != nil {
		m.updateChan <- Update{Type: UpdatedDevice, Device: dev, Changes: changes}
	}
}

// RemoveDevice removes a device based on the topic name
// Removing a device that does not exist is a no-op
func (m *Manager) removeDevice(topic string) {
	if !m.HasDevice(topic) {
		return
	}
	m.lock.Lock()
	dev := m.devices[topic]
	delete(m.devices, topic)
	if dev != nil {
		dev.stop()
	}
	m.lock.Unlock()
	if m.updateChan != nil {
		m.updateChan <- Update{Type: RemovedDevice, Device: dev}
	}
}

// UnreachableDevice marks the device specified by the topic
// as unreachable. This does not delete the device.
// Marking a non-existant device as unreachable is a no-op
func (m *Manager) unreachableDevice(topic string) {
	dev := m.Device(topic)
	if dev.Exists() {
		wasReachable := dev.IsReachable()
		dev.setReachability(false)
		if wasReachable && m.updateChan != nil {
			m.updateChan <- Update{
				Type:   UpdatedDevice,
				Device: dev,
				Changes: []*device.InfoUpdate{
					{Field: "reachable", Old: "1", New: "0"},
				},
			}
		}
	}
}

// Device returns a device associated with the topic
func (m *Manager) Device(topic string) Device {
	if m.HasDevice(topic) {
		m.lock.RLock()
		defer m.lock.RUnlock()
		return m.devices[topic]
	}
	err := errors.New("device not found")
	return &FakeDevice{Topic: topic, Err: err}
}

// Devices returns all devices known to the manager
func (m *Manager) Devices() []Device {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var devs []Device
	for _, dev := range m.devices {
		devs = append(devs, dev)
	}
	return devs
}

// WaitForDevice waits for a device on the specified topic to appear
// This blocks until either the device shows up, or the context is cancelled
func (m *Manager) WaitForDevice(ctx context.Context, topic string) Device {
	dev := m.Device(topic)
	if dev.Exists() {
		return dev
	}
	m.lock.Lock()

	ch, ok := m.waitingOn[topic]
	if !ok {
		ch = make(chan struct{})
		m.waitingOn[topic] = ch
	}

	m.lock.Unlock()

	select {
	case <-ch:
		return m.Device(topic)
	case <-ctx.Done():
		return &FakeDevice{Err: errors.New("context cancelled"), Topic: topic}
	}
}

// DeviceByType returns all devices known to the manager
// of type t.
func (m *Manager) DeviceByType(t string) []Device {
	devs := []Device{}
	for _, dev := range m.Devices() {
		if strings.ToLower(dev.Type()) == strings.ToLower(t) {
			devs = append(devs, dev)
		}
	}
	return devs
}
