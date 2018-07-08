package mqtt

import (
	"encoding/json"
	"github.com/goiiot/libmqtt"
	"github.com/hemtjanst/bibliotek/device"
	"github.com/hemtjanst/bibliotek/feature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
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
		deviceState:   make(chan *device.Info, 16),
		discoverDelay: 1 * time.Millisecond,
	}

	client.On("Connect", mock.Anything).Return()
	client.On("Publish", discoverTopic).Return()
	client.On("Subscribe", announceTopic+"/#").Return()
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

	assert.Equal(t, "I'm a teapot", dev.Name)
	time.Sleep(10 * time.Millisecond)
	client.AssertExpectations(t)
}
