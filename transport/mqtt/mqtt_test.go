package mqtt

import (
	"context"
	"errors"
	"sync"
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

func startMockClient(t *testing.T, client *MockMqttClient, expectErr error) (cl *mqtt, cancel func(), wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(context.Background())

	mqcl, err := New(ctx, &Config{Address: []string{"127.0.0.1:1883"}})
	assert.Nil(t, err)
	cl = mqcl.(*mqtt)
	cl.client = client
	cl.discoverDelay = 1 * time.Millisecond

	wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			ok, err := cl.Start()
			if !ok {
				return
			}
			assert.Equal(t, expectErr, err)
		}
	}()
	time.Sleep(10 * time.Millisecond)
	return cl, cancel, wg
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

	client.On("Connect", mock.Anything).Return()
	cl, cancel, wg := startMockClient(t, client, nil)

	client.AssertNotCalled(t, "Subscribe")
	client.On("Subscribe", "announce/#").Return()
	cl.DeviceState()
	client.AssertCalled(t, "Subscribe", "announce/#")
	client.AssertNotCalled(t, "Publish")
	client.On("Publish", "discover").Return()

	for i := 0; i < 2; i++ {
		time.Sleep(10 * time.Millisecond)
		client.AssertNumberOfCalls(t, "Publish", i+1)
		client.AssertNumberOfCalls(t, "Subscribe", i+1)
		client.AssertCalled(t, "Publish", "discover")
		client.TriggerReconnect(libmqtt.CodeSuccess, nil)
	}

	client.On("Destroy", true).Return()
	cancel()
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	client.AssertExpectations(t)
}

func TestClientError(t *testing.T) {
	expectErr := errors.New("test")

	client := &MockMqttClient{
		ConnectCode: libmqtt.CodeSuccess,
		ConnectErr:  expectErr,
	}

	client.On("Connect", mock.Anything).Return()
	client.On("Destroy", true).Return()
	_, cancel, wg := startMockClient(t, client, expectErr)

	cancel()
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	client.AssertExpectations(t)
}

func TestPublish(t *testing.T) {
	client := &MockMqttClient{
		ConnectCode: libmqtt.CodeSuccess,
	}

	client.On("Connect", mock.Anything).Return()
	cl, cancel, wg := startMockClient(t, client, nil)

	client.On("Publish", "harhartest").Return()
	client.On("Destroy", true).Return()
	cl.Publish("harhartest", []byte{}, true)
	cancel()
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	client.AssertExpectations(t)
}

func TestSubscribe(t *testing.T) {
	client := &MockMqttClient{
		ConnectCode: libmqtt.CodeSuccess,
	}

	client.On("Connect", mock.Anything).Return()
	cl, cancel, wg := startMockClient(t, client, nil)

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
	client.On("Destroy", true).Return()
	cancel()
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	client.AssertExpectations(t)
}
