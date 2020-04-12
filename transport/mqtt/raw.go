package mqtt

import (
	"github.com/goiiot/libmqtt"
	"strings"
)

func (m *mqtt) SubscribeRaw(topic string) chan *Packet {
	m.Lock()
	defer m.Unlock()
	if m.subRaw == nil || m.client == nil {
		// Probably not started
		return nil
	}
	c := make(chan *Packet, 5)

	if _, ok := m.subRaw[topic]; ok {
		m.subRaw[topic] = append(m.subRaw[topic], c)
		return c
	}
	m.subRaw[topic] = []chan *Packet{c}
	m.client.Subscribe(
		&libmqtt.Topic{Name: topic},
	)
	return c
}

func (m *mqtt) OnRaw(p *libmqtt.PublishPacket) {
	m.RLock()
	var chans []chan *Packet
	for k, v := range m.subRaw {
		if TopicTest(p.TopicName, k) {
			for _, ch := range v {
				chans = append(chans, ch)
			}
		}
	}
	m.RUnlock()
	for _, ch := range chans {
		func() {
			defer func() {
				_ = recover()
			}()
			ch <- (*Packet)(p)
		}()
	}

}

func TopicTest(topic, sub string) bool {
	subParts := strings.Split(sub, "/")
	topicParts := strings.Split(topic, "/")

	for i, p := range topicParts {
		if len(subParts) < i+1 {
			return false
		}
		sp := subParts[i]
		if p == sp || sp == "+" {
			continue
		}
		if sp == "#" && len(subParts) == i+1 {
			return true
		}
		return false
	}
	if len(subParts) > len(topicParts) {
		return false
	}
	return true
}
