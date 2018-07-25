package mqtt

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/goiiot/libmqtt"
	"github.com/hemtjanst/bibliotek/device"
)

type mqtt struct {
	deviceState   chan *device.Info
	client        mqttClient
	addr          string
	initCh        chan error
	sub           map[string][]chan []byte
	discoverSent  bool
	discoverDelay time.Duration
	sync.RWMutex
}

func New(ctx context.Context, addr string) (m *mqtt, err error) {
	m = &mqtt{
		addr:          addr,
		discoverDelay: 5 * time.Second,
	}
	opts := []libmqtt.Option{
		libmqtt.WithKeepalive(10, 1.2),
		libmqtt.WithLog(libmqtt.Silent),
		libmqtt.WithRouter(newRouter(m)),
		libmqtt.WithDialTimeout(5),
	}
	if addr != "" {
		opts = append(opts, libmqtt.WithServer(addr))
	}
	moreOpts, err := flagOpts()
	if err != nil {
		return nil, err
	}
	opts = append(opts, moreOpts...)

	client, err := libmqtt.NewClient(opts...)
	if err != nil {
		m = nil
		return
	}

	if err = m.init(ctx, client); err != nil {
		m = nil
	}

	return
}

func (m *mqtt) init(ctx context.Context, client mqttClient) (err error) {
	m.Lock()
	if m.client != nil {
		m.Unlock()
		return errors.New("already initialized")
	}
	m.initCh = make(chan error)
	m.sub = map[string][]chan []byte{}
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
		}()
	}

	return
}

func (m *mqtt) onConnect(server string, code byte, err error) {
	m.Lock()
	defer m.Unlock()
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

	for topic := range m.sub {
		m.client.Subscribe(
			&libmqtt.Topic{Name: topic},
		)
	}
}

func (m *mqtt) sendDiscover() {
	m.discoverSent = false
	m.client.Subscribe(&libmqtt.Topic{Name: announceTopic + "/#"})
	time.AfterFunc(m.discoverDelay, func() {
		m.Lock()
		defer m.Unlock()
		m.discoverSent = true
		m.client.Publish(&libmqtt.PublishPacket{
			TopicName: discoverTopic,
			IsRetain:  true,
			Payload:   []byte("1"),
		})
	})
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
	return make(chan struct{}, 5)
}
