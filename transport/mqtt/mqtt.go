package mqtt

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"fmt"
	"github.com/goiiot/libmqtt"
	"github.com/hemtjanst/bibliotek/device"
)

type mqtt struct {
	deviceState   chan *device.Info
	client        mqttClient
	addr          string
	initCh        chan error
	sub           map[string][]chan []byte
	discoverSub   []chan struct{}
	discoverSeen  bool
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
