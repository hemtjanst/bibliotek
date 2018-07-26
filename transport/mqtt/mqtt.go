package mqtt

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"fmt"
	"github.com/goiiot/libmqtt"
	"github.com/google/uuid"
	"github.com/hemtjanst/bibliotek/device"
	"os"
	"path"
)

type mqtt struct {
	deviceState   chan *device.Info
	client        mqttClient
	addr          string
	initCh        chan error
	sub           map[string][]chan []byte
	willID        string
	discoverSub   []chan struct{}
	discoverSeen  bool
	discoverSent  bool
	discoverDelay time.Duration
	sync.RWMutex
}

func New(ctx context.Context, addr string) (m *mqtt, err error) {
	var id string
	if len(os.Args) > 0 && len(os.Args[0]) > 0 {
		// Use executable name as first part of id
		id = path.Base(os.Args[0])
	} else {
		id = "htlib"
	}
	id = id + "-" + uuid.New().String()

	m = &mqtt{
		addr:          addr,
		discoverDelay: 5 * time.Second,
		willID:        id,
	}

	opts := []libmqtt.Option{
		libmqtt.WithKeepalive(10, 1.2),
		libmqtt.WithLog(libmqtt.Silent),
		libmqtt.WithRouter(newRouter(m)),
		libmqtt.WithDialTimeout(5),
		libmqtt.WithWill(leaveTopic, 0, false, []byte(m.willID)),
		libmqtt.WithClientID(id),
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

func (m *mqtt) LastWillID() string {
	return m.willID
}
