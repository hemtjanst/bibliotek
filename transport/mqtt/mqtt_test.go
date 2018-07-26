package mqtt

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/goiiot/libmqtt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		discoverDelay: 500 * time.Millisecond,
	}

	client.On("Connect", mock.Anything).Return()

	err := cl.init(ctx, client)
	assert.Nil(t, err)
	client.AssertNotCalled(t, "Subscribe")

	client.On("Subscribe", "announce/#").Return()
	cl.DeviceState()
	client.AssertCalled(t, "Subscribe", "announce/#")
	client.AssertNotCalled(t, "Publish")
	client.On("Publish", "discover").Return()

	for i := 0; i < 2; i++ {
		time.Sleep(1 * time.Second)
		client.AssertNumberOfCalls(t, "Publish", i+1)
		client.AssertNumberOfCalls(t, "Subscribe", i+1)
		client.AssertCalled(t, "Publish", "discover")
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
		discoverDelay: 500 * time.Millisecond,
	}

	client.On("Connect", mock.Anything).Return()
	client.On("Destroy", true).Return()
	err := cl.init(ctx, client)
	assert.Equal(t, expectErr, err)
	client.AssertExpectations(t)
}

func TestPublish(t *testing.T) {
	client := &MockMqttClient{
		ConnectCode: libmqtt.CodeSuccess,
	}

	cl := &mqtt{
		discoverDelay: 500 * time.Millisecond,
	}

	client.On("Connect", mock.Anything).Return()

	err := cl.init(context.Background(), client)
	assert.Nil(t, err)

	client.On("Publish", "harhartest").Return()
	cl.Publish("harhartest", []byte{}, true)
	client.AssertExpectations(t)
}

func TestSubscribe(t *testing.T) {
	client := &MockMqttClient{
		ConnectCode: libmqtt.CodeSuccess,
	}

	cl := &mqtt{
		discoverDelay: 500 * time.Millisecond,
	}

	client.On("Connect", mock.Anything).Return()

	err := cl.init(context.Background(), client)
	assert.Nil(t, err)

	client.On("Subscribe", "harhartest").Return()
	res1 := cl.Subscribe("harhartest")
	res2 := cl.Subscribe("harhartest")
	client.AssertNumberOfCalls(t, "Subscribe", 1)
	cl.OnFeature(&libmqtt.PublishPacket{
		TopicName: "harhartest",
		Payload:   []byte("test1"),
		IsRetain:  true,
	})

	select {
	case msg := <-res1:
		assert.Equal(t, []byte("test1"), msg)
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Reached timeout for res1")
	}
	select {
	case msg := <-res2:
		assert.Equal(t, []byte("test1"), msg)
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Reached timeout for res1")
	}

	client.TriggerReconnect(0x00, nil)
	client.On("Subscribe", "harhartest").Return()
	client.AssertNumberOfCalls(t, "Subscribe", 2)
	client.AssertExpectations(t)
}
