package mqtt

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/goiiot/libmqtt"
	"lib.hemtjan.st/device"
	"lib.hemtjan.st/feature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMessage(t *testing.T) {
	mockDev := func(t string, d *device.Info) *libmqtt.PublishPacket {
		data, _ := json.Marshal(d)
		return &libmqtt.PublishPacket{TopicName: t, Payload: data}
	}
	client := &MockMqttClient{
		ConnectCode: libmqtt.CodeSuccess,
	}

	cl := &mqtt{
		deviceState:   make(chan *device.State, 16),
		discoverDelay: 1 * time.Millisecond,
	}

	client.On("Connect", mock.Anything).Return()
	client.On("Publish", "discover").Return()
	client.On("Subscribe", "announce/#").Return()
	err := cl.init(nil, client)
	assert.Nil(t, err)

	cl.OnAnnounce(mockDev("announce/teapot", &device.Info{
		Name:         "I'm a teapot",
		Manufacturer: "foobar",
		Model:        "",
		LastWillID:   "xyz-123",
		Type:         "switch",
		Features: map[string]*feature.Info{
			"on": {},
		},
	}))

	dev := <-cl.deviceState

	assert.Equal(t, "I'm a teapot", dev.Device.Name)
	time.Sleep(10 * time.Millisecond)
	client.AssertExpectations(t)
}

func TestAnnounceLeave(t *testing.T) {
	mockDev := func(t string, d *device.Info) *libmqtt.PublishPacket {
		data, _ := json.Marshal(d)
		return &libmqtt.PublishPacket{TopicName: t, Payload: data}
	}
	cl := &mqtt{
		announceTopic: "announce",
		deviceState:   make(chan *device.State, 16),
		discoverDelay: 1 * time.Millisecond,
	}
	cl.OnAnnounce(mockDev("announce/teapot", &device.Info{
		Name:       "I'm a teapot",
		LastWillID: "xyz-123",
		Features:   map[string]*feature.Info{},
	}))
	dev := <-cl.deviceState
	assert.Equal(t, "xyz-123", dev.Device.LastWillID)
	assert.Equal(t, map[string][]string{"xyz-123": {"teapot"}}, cl.willMap)
	cl.OnLeave(&libmqtt.PublishPacket{TopicName: "leave", Payload: []byte("xyz-123")})
	dev = <-cl.deviceState
	assert.Equal(t, &device.State{Action: device.LeaveAction, Topic: "teapot"}, dev)
	assert.Equal(t, map[string][]string{}, cl.willMap)

}

func TestAnnounceNewWillID(t *testing.T) {
	mockDev := func(t string, d *device.Info) *libmqtt.PublishPacket {
		data, _ := json.Marshal(d)
		return &libmqtt.PublishPacket{TopicName: t, Payload: data}
	}
	cl := &mqtt{
		announceTopic: "announce",
		deviceState:   make(chan *device.State, 16),
		discoverDelay: 1 * time.Millisecond,
	}
	cl.OnAnnounce(mockDev("announce/teapot1", &device.Info{
		Name:       "I'm a teapot",
		LastWillID: "xyz-123",
		Features:   map[string]*feature.Info{},
	}))
	cl.OnAnnounce(mockDev("announce/teapot2", &device.Info{
		Name:       "I'm a teapot",
		LastWillID: "xyz-123",
		Features:   map[string]*feature.Info{},
	}))
	dev := <-cl.deviceState
	assert.Equal(t, "xyz-123", dev.Device.LastWillID)
	assert.Equal(t, "teapot1", dev.Topic)
	dev = <-cl.deviceState
	assert.Equal(t, "xyz-123", dev.Device.LastWillID)
	assert.Equal(t, "teapot2", dev.Topic)
	assert.Equal(t, map[string][]string{"xyz-123": {"teapot1", "teapot2"}}, cl.willMap)

	cl.OnAnnounce(mockDev("announce/teapot1", &device.Info{
		Name:       "I'm a teapot",
		LastWillID: "abc-123",
		Features:   map[string]*feature.Info{},
	}))

	dev = <-cl.deviceState
	assert.Equal(t, "abc-123", dev.Device.LastWillID)
	assert.Equal(t, "teapot1", dev.Topic)
	assert.Equal(t, map[string][]string{"abc-123": {"teapot1"}, "xyz-123": {"teapot2"}}, cl.willMap)

	cl.OnAnnounce(mockDev("announce/teapot2", &device.Info{
		Name:       "I'm a teapot",
		LastWillID: "abc-123",
		Features:   map[string]*feature.Info{},
	}))

	dev = <-cl.deviceState
	assert.Equal(t, "abc-123", dev.Device.LastWillID)
	assert.Equal(t, "teapot2", dev.Topic)
	assert.Equal(t, map[string][]string{"abc-123": {"teapot1", "teapot2"}}, cl.willMap)

	cl.OnLeave(&libmqtt.PublishPacket{TopicName: "leave", Payload: []byte("xyz-123")})
	assert.Equal(t, map[string][]string{"abc-123": {"teapot1", "teapot2"}}, cl.willMap)
}

func TestAnnounceDelete(t *testing.T) {
	mockDev := func(t string, d *device.Info) *libmqtt.PublishPacket {
		data, _ := json.Marshal(d)
		return &libmqtt.PublishPacket{TopicName: t, Payload: data}
	}
	cl := &mqtt{
		announceTopic: "announce",
		deviceState:   make(chan *device.State, 16),
		discoverDelay: 1 * time.Millisecond,
	}
	cl.OnAnnounce(mockDev("announce/teapot1", &device.Info{
		Name:       "I'm a teapot",
		LastWillID: "xyz-123",
		Features:   map[string]*feature.Info{},
	}))
	dev := <-cl.deviceState
	assert.Equal(t, "xyz-123", dev.Device.LastWillID)
	assert.Equal(t, "teapot1", dev.Topic)
	assert.Equal(t, map[string][]string{"xyz-123": {"teapot1"}}, cl.willMap)

	cl.OnAnnounce(&libmqtt.PublishPacket{TopicName: "announce/teapot1", Payload: []byte{}, IsRetain: true})
	dev = <-cl.deviceState
	assert.Equal(t, device.DeleteAction, dev.Action)
	assert.Equal(t, "teapot1", dev.Topic)
	assert.Equal(t, map[string][]string{}, cl.willMap)

}
