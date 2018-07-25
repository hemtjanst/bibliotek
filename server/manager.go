package server

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/hemtjanst/bibliotek/device"
)

type Transport interface {
	// DeviceState receives a device.Info on announce, leave or update
	DeviceState() chan *device.Info
	device.Transport
}

// Manager holds all the devices and deals with updating device
// state as it receives it from the transport
type Manager struct {
	transport Transport
	devices   map[string]device.Device
	waitingOn map[string]chan struct{}
	sync.RWMutex
}

// New creates a new Manager
func New(t Transport) *Manager {
	return &Manager{
		transport: t,
		devices:   map[string]device.Device{},
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
			log.Printf("Device State: %v", d)
			if !d.Reachable {
				continue
			}
			if !m.HasDevice(d.Topic) {
				m.AddDevice(d)
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
	dev, err := device.NewServer(d, m.transport)
	if err != nil {
		log.Printf("Failed to create device: %v", err)
		return
	}
	m.Lock()
	defer m.Unlock()
	m.devices[d.Topic] = dev
	if ch, ok := m.waitingOn[d.Topic]; ok {
		close(ch)
	}
}

func (m *Manager) Device(topic string) device.Device {
	if m.HasDevice(topic) {
		m.RLock()
		defer m.RUnlock()
		return m.devices[topic]
	}
	err := errors.New("device not found")
	return &device.Fake{Topic: topic, Err: err}
}

func (m *Manager) Devices() []device.Device {
	m.RLock()
	defer m.RUnlock()
	var devs []device.Device
	for _, dev := range m.devices {
		devs = append(devs, dev)
	}
	return devs
}

func (m *Manager) WaitForDevice(topic string, ctx context.Context) device.Device {
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
		return &device.Fake{Err: errors.New("context cancelled"), Topic: topic}
	}
}
