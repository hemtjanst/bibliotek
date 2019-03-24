package mqtt

import (
	"github.com/goiiot/libmqtt"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockMessageHandler struct{ mock.Mock }

func (h *MockMessageHandler) OnRaw(p *libmqtt.PublishPacket)      {}
func (h *MockMessageHandler) OnDiscover(p *libmqtt.PublishPacket) { h.Called(p.TopicName) }
func (h *MockMessageHandler) OnAnnounce(p *libmqtt.PublishPacket) { h.Called(p.TopicName) }
func (h *MockMessageHandler) OnLeave(p *libmqtt.PublishPacket)    { h.Called(p.TopicName) }
func (h *MockMessageHandler) OnFeature(p *libmqtt.PublishPacket)  { h.Called(p.TopicName) }
func (m *MockMessageHandler) TopicName(t EventType) string {
	switch t {
	case TypeAnnounce:
		return "announce"
	case TypeDiscover:
		return "discover"
	case TypeLeave:
		return "leave"
	default:
		return ""
	}
}

func TestRouter(t *testing.T) {
	mockPacket := func(topic string) *libmqtt.PublishPacket {
		return &libmqtt.PublishPacket{
			TopicName: topic,
			Payload:   []byte("foo"),
		}
	}
	handler := &MockMessageHandler{}
	r := newRouter(handler)
	handler.On("OnDiscover", "discover").Return()
	handler.On("OnAnnounce", "announce/teapot").Return()
	handler.On("OnLeave", "leave").Return()
	handler.On("OnFeature", "anything").Return()
	r.Dispatch(mockPacket("announce/teapot"))
	r.Dispatch(mockPacket("discover"))
	r.Dispatch(mockPacket("leave"))
	r.Dispatch(mockPacket("anything"))

	handler.AssertExpectations(t)

}
