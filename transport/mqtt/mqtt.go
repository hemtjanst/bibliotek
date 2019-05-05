package mqtt

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/goiiot/libmqtt"
	"lib.hemtjan.st/device"
)

type Packet libmqtt.PublishPacket

type mqtt struct {
	deviceState   chan *device.State
	client        mqttClient
	addr          string
	initCh        chan error
	sub           map[string][]chan []byte
	subRaw        map[string][]chan *Packet
	willMap       map[string][]string
	willID        string
	discoverSub   []chan struct{}
	discoverSeen  bool
	discoverSent  bool
	discoverDelay time.Duration
	announceTopic string
	discoverTopic string
	leaveTopic    string
	sync.RWMutex
}

func New(ctx context.Context, c *Config) (m MQTT, err error) {
	if c == nil {
		return nil, ErrNoConfig
	}
	if err = c.check(); err != nil {
		return
	}
	mq := &mqtt{
		discoverDelay: c.DiscoverDelay,
		willID:        c.ClientID,
		announceTopic: c.AnnounceTopic,
		discoverTopic: c.DiscoverTopic,
		leaveTopic:    c.LeaveTopic,
		willMap:       map[string][]string{},
	}
	opts := []libmqtt.Option{
		libmqtt.WithRouter(newRouter(mq)),
	}
	opts = append(opts, c.opts()...)

	client, err := libmqtt.NewClient(opts...)
	if err != nil {
		m = nil
		return
	}

	if err = mq.init(ctx, client); err != nil {
		m = nil
	}
	m = mq
	return
}

func (m *mqtt) init(ctx context.Context, client mqttClient) (err error) {
	m.Lock()
	if m.client != nil {
		m.Unlock()
		return errors.New("already initialized")
	}
	if m.announceTopic == "" {
		m.announceTopic = "announce"
	}
	if m.discoverTopic == "" {
		m.discoverTopic = "discover"
	}
	if m.leaveTopic == "" {
		m.leaveTopic = "leave"
	}
	if m.discoverDelay == 0 {
		m.discoverDelay = 5 * time.Second
	}

	m.initCh = make(chan error)
	m.sub = map[string][]chan []byte{}
	m.subRaw = map[string][]chan *Packet{}
	m.client = client
	m.Unlock()
	m.client.Connect(m.onConnect)
	err, _ = <-m.initCh
	m.initCh = nil

	if err != nil {
		m.client.Destroy(true)
		return
	}

	if ctx != nil {
		go func() {
			<-ctx.Done()
			m.client.Destroy(false)
			m.Lock()
			dsub := m.discoverSub
			m.discoverSub = []chan struct{}{}
			subs := m.sub
			m.sub = map[string][]chan []byte{}
			stateCh := m.deviceState
			m.deviceState = nil
			m.Unlock()

			if stateCh != nil {
				close(stateCh)
			}

			for _, ch := range dsub {
				close(ch)
			}

			if m.subRaw != nil {
				for _, v := range m.subRaw {
					for _, vv := range v {
						close(vv)
					}
				}
				m.subRaw = map[string][]chan *Packet{}
			}

			if subs != nil {
				for _, chans := range subs {
					for _, ch := range chans {
						close(ch)
					}
				}
			}
		}()
	}

	return
}

func (m *mqtt) onConnect(server string, code byte, err error) {
	m.Lock()
	defer m.Unlock()

	if code != libmqtt.CodeSuccess && err == nil {
		err = fmt.Errorf("error code %d", int(code))
	}

	if m.initCh != nil {
		if err != nil {
			m.initCh <- err
		} else {
			close(m.initCh)
		}
	}

	if err != nil {
		log.Printf("MQTT Connect Error: %s (%x) %v", server, code, err)
		return
	}

	if m.deviceState != nil {
		m.sendDiscover()
	}

	seen := map[string]bool{}
	for topic := range m.sub {
		seen[topic] = true
		m.client.Subscribe(
			&libmqtt.Topic{Name: topic},
		)
	}
	for topic := range m.subRaw {
		if _, ok := seen[topic]; ok {
			continue
		}
		m.client.Subscribe(
			&libmqtt.Topic{Name: topic},
		)
	}
}

func (m *mqtt) sendDiscover() {
	m.discoverSent = false
	m.client.Subscribe(&libmqtt.Topic{Name: m.announceTopic + "/#"})
	time.AfterFunc(m.discoverDelay, func() {
		m.Lock()
		defer m.Unlock()
		m.discoverSent = true
		m.client.Publish(&libmqtt.PublishPacket{
			TopicName: m.discoverTopic,
			IsRetain:  true,
			Payload:   []byte("1"),
		})
	})
}

func (m *mqtt) LastWillID() string {
	return m.willID
}

func (m *mqtt) updateWills(devTopic string, newWillID string) {
	m.Lock()
	defer m.Unlock()
	if m.willMap == nil {
		m.willMap = map[string][]string{}
	}
	found := false
outer:
	for k, v := range m.willMap {
		var mm []string
		for _, vv := range v {
			if vv == devTopic {
				if k == newWillID {
					found = true
					continue outer
				}
				continue
			}
			mm = append(mm, vv)
		}
		if k == newWillID {
			found = true
			mm = append(mm, devTopic)
		}
		if len(mm) == 0 {
			delete(m.willMap, k)
		} else if len(mm) != len(v) {
			m.willMap[k] = mm
		}
	}
	if !found && len(newWillID) > 0 {
		m.willMap[newWillID] = []string{devTopic}
	}
}
