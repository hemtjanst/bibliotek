package mqtt

import (
	"encoding/json"
	"github.com/goiiot/libmqtt"
	"github.com/hemtjanst/bibliotek/device"
	"testing"
	"time"
)

func TestMessage(t *testing.T) {
	mockDev := func(t string, d *device.DeviceInfo) *libmqtt.PublishPacket {
		data, _ := json.Marshal(d)
		return &libmqtt.PublishPacket{TopicName: t, Payload: data}
	}
	mock := &MockMqttClient{
		DestroyChan: make(chan bool, 5),
		ConnectCode: libmqtt.CodeSuccess,
	}

	cl := &mqtt{
		deviceState:   make(chan *device.DeviceInfo, 16),
		discoverDelay: 1 * time.Second,
	}
	err := cl.init(nil, mock)
	if err != nil {
		t.Errorf("Got error from init: %v", err)
	}

	cl.OnAnnounce(mockDev("announce/teapot", &device.DeviceInfo{
		Name:         "I'm a teapot",
		Manufacturer: "foobar",
		Model:        "",
		LastWillID:   "xyz-123",
		Type:         "switch",
		Features: map[string]*device.FeatureInfo{
			"on": {},
		},
	}))

	dev := <-cl.deviceState

	if dev.Reachable {
		t.Errorf("Device shouldn't be reachable")
	}
	if dev.Name != "I'm a teapot" {
		t.Errorf("Name is '%s'", dev.Name)
	}

}
