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
	sendDiscover := false
	if m.deviceState == nil {
		m.deviceState = make(chan *device.State, 10)
		sendDiscover = true
	}
	ds := m.deviceState
	m.Unlock()
	if sendDiscover {
		m.sendDiscover()
	}
	return ds
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
	if v, ok := m.sub[topic]; ok {
		defer func() {
			for _, ch := range v {
				close(ch)
			}
		}()
		delete(m.sub, topic)
		found = true
	}
	if m.subRaw != nil {
		if v, ok := m.subRaw[topic]; ok {
			defer func() {
				for _, ch := range v {
					close(ch)
				}
			}()
			delete(m.sub, topic)
			found = true
		}
	}
	m.Unlock()
	if found {
		m.client.UnSubscribe(topic)
	}
	return
}

func (m *mqtt) Resubscribe(oldTopic, newTopic string) bool {
	m.Lock()
	keep := false
	subNew := false
	if m.subRaw != nil {
		_, keep = m.subRaw[oldTopic]
	}
	if v, ok := m.sub[oldTopic]; ok {
		if _, ok := m.sub[newTopic]; !ok {
			m.sub[newTopic] = v
			subNew = true
		} else {
			m.sub[newTopic] = append(m.sub[newTopic], v...)
		}
		delete(m.sub, oldTopic)
		m.Unlock()
		if !keep {
			m.client.UnSubscribe(oldTopic)
		}
		if subNew {
			m.client.Subscribe(
				&libmqtt.Topic{Name: newTopic},
			)
		}
		return true
	} else {
		m.Unlock()
	}
	return false
}

func (m *mqtt) Subscribe(topic string) chan []byte {
	m.Lock()
	c := make(chan []byte, 5)

	if _, ok := m.sub[topic]; ok {
		m.sub[topic] = append(m.sub[topic], c)
		m.Unlock()
		return c
	}
	m.sub[topic] = []chan []byte{c}
	m.Unlock()
	m.client.Subscribe(
		&libmqtt.Topic{Name: topic},
	)
	return c
}

func (m *mqtt) Discover() chan struct{} {
	m.Lock()
	ch := make(chan struct{}, 5)
	m.discoverSub = append(m.discoverSub, ch)
	discLen := len(m.discoverSub)
	if m.discoverSeen {
		ch <- struct{}{}
	}
	m.Unlock()
	if discLen == 1 {
		m.client.Subscribe(&libmqtt.Topic{Name: m.discoverTopic})
	}
	return ch
}
