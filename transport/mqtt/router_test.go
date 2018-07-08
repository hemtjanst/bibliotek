package mqtt

import (
	"encoding/json"
	"github.com/goiiot/libmqtt"
	"github.com/hemtjanst/bibliotek/device"
	"testing"
)

type MockMessageHandler struct {
	Discover []*libmqtt.PublishPacket
	Announce []*libmqtt.PublishPacket
	Leave    []*libmqtt.PublishPacket
	Feature  []*libmqtt.PublishPacket
}

func (h *MockMessageHandler) OnDiscover(p *libmqtt.PublishPacket) { h.Discover = append(h.Discover, p) }
func (h *MockMessageHandler) OnAnnounce(p *libmqtt.PublishPacket) { h.Announce = append(h.Announce, p) }
func (h *MockMessageHandler) OnLeave(p *libmqtt.PublishPacket)    { h.Leave = append(h.Leave, p) }
func (h *MockMessageHandler) OnFeature(p *libmqtt.PublishPacket)  { h.Feature = append(h.Feature, p) }
func (h *MockMessageHandler) CheckCount(d, a, l, f int) bool {
	return len(h.Discover) == d && len(h.Announce) == a && len(h.Leave) == l && len(h.Feature) == f
}

func TestRouter(t *testing.T) {
	mockPacket := func(topic string, payload []byte) *libmqtt.PublishPacket {
		return &libmqtt.PublishPacket{
			TopicName: topic,
			Payload:   []byte(payload),
		}
	}
	mock := &MockMessageHandler{}
	r := newRouter(mock)
	dev := &device.DeviceInfo{Name: "teapot", Manufacturer: "foobar"}
	data, _ := json.Marshal(dev)
	r.Dispatch(mockPacket("announce/teapot", data))

	if !mock.CheckCount(0, 1, 0, 0) {
		t.Errorf("Invalid message count (1)")
	}

	r.Dispatch(mockPacket("discover", []byte("1")))
	if !mock.CheckCount(1, 1, 0, 0) {
		t.Errorf("Invalid message count (2)")
	}

	r.Dispatch(mockPacket("leave", []byte("teapot")))
	if !mock.CheckCount(1, 1, 1, 0) {
		t.Errorf("Invalid message count (3)")
	}

	r.Dispatch(mockPacket("anything", []byte("foo")))
	if !mock.CheckCount(1, 1, 1, 1) {
		t.Errorf("Invalid message count (4)")
	}

}
