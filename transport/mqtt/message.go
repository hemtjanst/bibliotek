package mqtt

import (
	"encoding/json"
	"github.com/goiiot/libmqtt"
	"github.com/hemtjanst/bibliotek/device"
	"log"
)

type MessageHandler interface {
	OnAnnounce(p *libmqtt.PublishPacket)
	OnLeave(p *libmqtt.PublishPacket)
	OnDiscover(p *libmqtt.PublishPacket)
	OnFeature(p *libmqtt.PublishPacket)
}

const (
	announceTopic = "announce"
	leaveTopic    = "leave"
	discoverTopic = "discover"
)

func (m *mqtt) OnAnnounce(p *libmqtt.PublishPacket) {
	if m.deviceState == nil {
		return
	}
	dev := &device.DeviceInfo{}
	err := json.Unmarshal(p.Payload, dev)
	if dev.Topic == "" {
		dev.Topic = p.TopicName[len(announceTopic)+1:]
	}
	dev.Reachable = m.discoverSent
	if err != nil {
		log.Printf("Error in json: %v (packet: %s)", err, string(p.Payload))
	}
	m.deviceState <- dev
}

func (m *mqtt) OnLeave(p *libmqtt.PublishPacket) {

}

func (m *mqtt) OnDiscover(p *libmqtt.PublishPacket) {

}

func (m *mqtt) OnFeature(p *libmqtt.PublishPacket) {

}

// DeviceState returns a channel which publishes information about new and changed devices
func (m *mqtt) DeviceState() chan *device.DeviceInfo {
	if m.deviceState == nil {
		m.deviceState = make(chan *device.DeviceInfo)
		m.sendDiscover()
	}
	return m.deviceState
}
