package mqtt

import (
	"context"
	"errors"
	"github.com/goiiot/libmqtt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockMqttClient struct {
	mock.Mock
	ConnectErr  error
	ConnectCode byte
	connHandler libmqtt.ConnHandler
}

func (m *MockMqttClient) TriggerReconnect(code byte, err error) {
	m.connHandler("127.0.0.1:1883", code, err)
}

func (m *MockMqttClient) Connect(handler libmqtt.ConnHandler) {
	m.Called(handler)
	m.connHandler = handler
	go func() {
		handler("127.0.0.1:1883", m.ConnectCode, m.ConnectErr)
	}()
}
func (m *MockMqttClient) Publish(p ...*libmqtt.PublishPacket) {
	for _, msg := range p {
		if msg.TopicName == "devnull" {
			continue
		}
		m.Called(msg.TopicName)
	}

}
func (m *MockMqttClient) Subscribe(topic ...*libmqtt.Topic) {
	m.Called(topic[0].Name)
}

func (m *MockMqttClient) UnSubscribe(topic ...string) {
	m.Called(topic[0])
}

func (m *MockMqttClient) Destroy(force bool) {
	m.Called(force)
}

func TestClient(t *testing.T) {
	client := &MockMqttClient{
		ConnectCode: libmqtt.CodeSuccess,
	}
	ctx, cancel := context.WithCancel(context.Background())

	cl := &mqtt{
		addr:          "127.0.0.1:1883",
		discoverDelay: 500 * time.Millisecond,
	}

	client.On("Connect", mock.Anything).Return()

	err := cl.init(ctx, client)
	assert.Nil(t, err)
	client.AssertNotCalled(t, "Subscribe")

	client.On("Subscribe", announceTopic+"/#").Return()
	cl.DeviceState()
	client.AssertCalled(t, "Subscribe", announceTopic+"/#")
	client.AssertNotCalled(t, "Publish")
	client.On("Publish", discoverTopic).Return()

	for i := 0; i < 2; i++ {
		time.Sleep(1 * time.Second)
		client.AssertNumberOfCalls(t, "Publish", i+1)
		client.AssertNumberOfCalls(t, "Subscribe", i+1)
		client.AssertCalled(t, "Publish", discoverTopic)
		client.TriggerReconnect(libmqtt.CodeSuccess, nil)
	}

	client.On("Destroy", false).Return()
	cancel()
	time.Sleep(5 * time.Millisecond)
	client.AssertExpectations(t)
}

func TestClientError(t *testing.T) {
	expectErr := errors.New("test")

	client := &MockMqttClient{
		ConnectCode: libmqtt.CodeSuccess,
		ConnectErr:  expectErr,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cl := &mqtt{
		addr:          "127.0.0.1:1883",
		discoverDelay: 500 * time.Millisecond,
	}

	client.On("Connect", mock.Anything).Return()
	client.On("Destroy", true).Return()
	err := cl.init(ctx, client)
	assert.Equal(t, expectErr, err)
	client.AssertExpectations(t)
}
