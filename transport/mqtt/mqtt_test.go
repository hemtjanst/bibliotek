package mqtt

import (
	"context"
	"errors"
	"github.com/goiiot/libmqtt"
	"testing"
	"time"
)

type MockMqttClient struct {
	ConnectErr    error
	ConnectCode   byte
	DestroyChan   chan bool
	Subscriptions []*libmqtt.Topic
	Published     []*libmqtt.PublishPacket
	connHandler   libmqtt.ConnHandler
}

func (m *MockMqttClient) TriggerReconnect(code byte, err error) {
	m.Subscriptions = []*libmqtt.Topic{}
	m.Published = []*libmqtt.PublishPacket{}
	m.connHandler("127.0.0.1:1883", code, err)
}

func (m *MockMqttClient) Connect(handler libmqtt.ConnHandler) {
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
		m.Published = append(m.Published, msg)
	}

}
func (m *MockMqttClient) Subscribe(topic ...*libmqtt.Topic) {
	m.Subscriptions = append(m.Subscriptions, topic...)
}
func (m *MockMqttClient) UnSubscribe(topic ...string) {
	var save []*libmqtt.Topic
Next:
	for _, v := range m.Subscriptions {
		for _, t := range topic {
			if v.Name == t {
				continue Next
			}
		}
		save = append(save, v)
	}
	m.Subscriptions = save
}
func (m *MockMqttClient) Destroy(force bool) {
	if m.DestroyChan != nil {
		m.DestroyChan <- force
	}
}

func TestClient(t *testing.T) {
	mock := &MockMqttClient{
		DestroyChan: make(chan bool),
		ConnectCode: libmqtt.CodeSuccess,
	}
	ctx, cancel := context.WithCancel(context.Background())

	cl := &mqtt{
		addr:          "127.0.0.1:1883",
		discoverDelay: 500 * time.Millisecond,
	}
	err := cl.init(ctx, mock)
	if err != nil {
		t.Errorf("Client init returned %v", err)
	}

	if len(mock.Subscriptions) > 0 {
		t.Errorf("Client shouldn't have any subscriptions")
	}

	cl.DeviceState()
	for i := 0; i < 2; i++ {
		if len(mock.Subscriptions) != 1 {
			t.Errorf("[%d] Client should have one subscription", i)
		} else if mock.Subscriptions[0].Name != announceTopic+"/#" {
			t.Errorf("[%d] Client should subscribe to '%s/#'", i, announceTopic)
		}

		if len(mock.Published) > 0 {
			t.Errorf("[%d] No messages should've been published", i)
		}

		time.Sleep(1 * time.Second)
		if len(mock.Published) != 1 {
			t.Errorf("[%d] Exactly one messaged should've been published (have %d)", i, len(mock.Published))
		} else if mock.Published[0].TopicName != discoverTopic {
			t.Errorf("[%d] Discover message sent to wrong topic (%s)", i, mock.Published[0].TopicName)
		}
		mock.TriggerReconnect(libmqtt.CodeSuccess, nil)
	}

	cancel()
	select {
	case forceFlag := <-mock.DestroyChan:
		if forceFlag {
			t.Errorf("destroy() was called with force=true")
		}
	case <-time.After(3 * time.Second):
		t.Errorf("cancel() didn't call destroy()")
	}

}

func TestClientError(t *testing.T) {
	expectErr := errors.New("test")

	mock := &MockMqttClient{
		DestroyChan: make(chan bool, 5),
		ConnectCode: libmqtt.CodeSuccess,
		ConnectErr:  expectErr,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cl := &mqtt{
		addr:          "127.0.0.1:1883",
		discoverDelay: 500 * time.Millisecond,
	}

	err := cl.init(ctx, mock)
	if err != expectErr {
		t.Errorf("Wrong error, expected '%v', got '%v'", expectErr, err)
	}
	force := <-mock.DestroyChan
	if !force {
		t.Errorf("Force destroy not used on initial error")
	}

}
