package mqtt

import (
	"encoding/json"
	"log"

	"github.com/goiiot/libmqtt"
	"github.com/hemtjanst/bibliotek/device"
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
	m.RLock()
	if m.deviceState == nil {
		m.RUnlock()
		return
	}
	reachable := m.discoverSent
	m.RUnlock()
	dev := &device.Info{}
	err := json.Unmarshal(p.Payload, dev)
	if dev.Topic == "" {
		dev.Topic = p.TopicName[len(announceTopic)+1:]
	}
	dev.Reachable = reachable
	if err != nil {
		log.Printf("Error in json: %v (packet: %s)", err, string(p.Payload))
		return
	}
	m.deviceState <- dev
}

func (m *mqtt) OnLeave(p *libmqtt.PublishPacket) {

}

func (m *mqtt) OnDiscover(p *libmqtt.PublishPacket) {
	m.RLock()
	chans := m.discoverSub
	seen := m.discoverSeen
	m.RUnlock()
	if !seen {
		m.Lock()
		m.discoverSeen = true
		m.Unlock()
	}
	for _, ch := range chans {
		ch <- struct{}{}
	}
}

func (m *mqtt) OnFeature(p *libmqtt.PublishPacket) {
	m.RLock()
	s, ok := m.sub[p.TopicName]
	m.RUnlock()
	if ok {
		for _, f := range s {
			f <- p.Payload
		}
	}
}

// DeviceState returns a channel which publishes information about new and changed devices
func (m *mqtt) DeviceState() chan *device.Info {
	m.Lock()
	defer m.Unlock()
	if m.deviceState == nil {
		m.deviceState = make(chan *device.Info, 10)
		m.sendDiscover()
	}
	return m.deviceState
}

func (m *mqtt) PublishMeta(topic string, payload []byte) {
	m.client.Publish(
		&libmqtt.PublishPacket{
			TopicName: announceTopic + "/" + topic,
			Payload:   payload,
			IsRetain:  true,
		},
	)
}

func (m *mqtt) Publish(topic string, payload []byte, retain bool) {
	m.client.Publish(
		&libmqtt.PublishPacket{
			TopicName: topic,
			Payload:   payload,
			IsRetain:  retain,
		},
	)
}
func (m *mqtt) Subscribe(topic string) chan []byte {
	m.Lock()
	defer m.Unlock()
	c := make(chan []byte, 5)

	if _, ok := m.sub[topic]; ok {
		m.sub[topic] = append(m.sub[topic], c)
		return c
	}
	m.sub[topic] = []chan []byte{c}
	m.client.Subscribe(
		&libmqtt.Topic{Name: topic},
	)
	return c
}

func (m *mqtt) Discover() chan struct{} {
	m.Lock()
	defer m.Unlock()
	ch := make(chan struct{}, 5)
	m.discoverSub = append(m.discoverSub, ch)
	if m.discoverSeen {
		ch <- struct{}{}
	}
	if len(m.discoverSub) == 1 {
		m.client.Subscribe(&libmqtt.Topic{Name: discoverTopic})
	}
	return ch
}
