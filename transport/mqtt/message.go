package mqtt

import (
	"encoding/json"
	"log"

	"github.com/goiiot/libmqtt"
	"lib.hemtjan.st/device"
)

type EventType int

const (
	TypeAnnounce EventType = iota
	TypeDiscover
	TypeLeave
)

type MessageHandler interface {
	OnRaw(p *libmqtt.PublishPacket)
	OnAnnounce(p *libmqtt.PublishPacket)
	OnLeave(p *libmqtt.PublishPacket)
	OnDiscover(p *libmqtt.PublishPacket)
	OnFeature(p *libmqtt.PublishPacket)
	TopicName(t EventType) string
}

func (m *mqtt) OnAnnounce(p *libmqtt.PublishPacket) {
	m.RLock()
	if m.deviceState == nil {
		m.RUnlock()
		return
	}
	reachable := m.discoverSent
	m.RUnlock()
	devTopic := p.TopicName[len(m.announceTopic)+1:]

	if len(p.Payload) == 0 {
		m.updateWills(devTopic, "")
		m.deviceState <- &device.State{
			Topic:  devTopic,
			Action: device.DeleteAction,
		}
		return
	}

	dev := &device.Info{}
	err := json.Unmarshal(p.Payload, dev)
	if dev.Topic == "" {
		dev.Topic = devTopic
	}
	dev.Reachable = reachable
	if err != nil {
		log.Printf("Error in json: %v (packet: %s)", err, string(p.Payload))
		return
	}
	m.updateWills(devTopic, dev.LastWillID)
	m.deviceState <- &device.State{
		Device: dev,
		Action: device.UpdateAction,
		Topic:  devTopic,
	}
}

func (m *mqtt) OnLeave(p *libmqtt.PublishPacket) {
	willID := string(p.Payload)
	if willID == "" {
		return
	}
	m.RLock()
	if m.willMap == nil {
		m.RUnlock()
		return
	}
	var devices []string
	var ok bool

	if devices, ok = m.willMap[string(p.Payload)]; !ok {
		m.RUnlock()
		return
	}
	m.RUnlock()
	m.Lock()
	delete(m.willMap, willID)
	m.Unlock()

	for _, d := range devices {
		m.deviceState <- &device.State{
			Topic:  d,
			Action: device.LeaveAction,
		}
	}
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

func (m *mqtt) TopicName(t EventType) string {
	switch t {
	case TypeAnnounce:
		return m.announceTopic
	case TypeDiscover:
		return m.discoverTopic
	case TypeLeave:
		return m.leaveTopic
	default:
		return ""
	}
}

// DeviceState returns a channel which publishes information about new and changed devices
func (m *mqtt) DeviceState() chan *device.State {
	m.Lock()
	defer m.Unlock()
	if m.deviceState == nil {
		m.deviceState = make(chan *device.State, 10)
		m.sendDiscover()
	}
	return m.deviceState
}

func (m *mqtt) PublishMeta(topic string, payload []byte) {
	m.client.Publish(
		&libmqtt.PublishPacket{
			TopicName: m.announceTopic + "/" + topic,
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

// TODO: Make UnsubscribeRaw()?

func (m *mqtt) Unsubscribe(topic string) (found bool) {
	m.Lock()
	defer m.Unlock()
	if v, ok := m.sub[topic]; ok {
		for _, ch := range v {
			close(ch)
		}
		delete(m.sub, topic)
		m.client.UnSubscribe(topic)
		found = true
	}
	if m.subRaw == nil {
		return
	}
	if v, ok := m.subRaw[topic]; ok {
		for _, ch := range v {
			close(ch)
		}
		delete(m.sub, topic)
		if !found {
			m.client.UnSubscribe(topic)
			found = true
		}
	}
	return
}

func (m *mqtt) Resubscribe(oldTopic, newTopic string) bool {
	m.Lock()
	defer m.Unlock()
	keep := false
	if m.subRaw != nil {
		_, keep = m.subRaw[oldTopic]
	}
	if v, ok := m.sub[oldTopic]; ok {
		if _, ok := m.sub[newTopic]; !ok {
			m.sub[newTopic] = v
			m.client.Subscribe(
				&libmqtt.Topic{Name: newTopic},
			)
		} else {
			m.sub[newTopic] = append(m.sub[newTopic], v...)
		}
		delete(m.sub, oldTopic)
		if !keep {
			m.client.UnSubscribe(oldTopic)
		}
		return true
	}
	return false
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
		m.client.Subscribe(&libmqtt.Topic{Name: m.discoverTopic})
	}
	return ch
}
