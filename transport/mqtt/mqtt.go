package mqtt

import (
	"context"
	"github.com/goiiot/libmqtt"
	"github.com/hemtjanst/bibliotek/device"
	"log"
	"time"
)

type mqtt struct {
	deviceState   chan *device.DeviceInfo
	client        mqttClient
	addr          string
	initCh        chan error
	sub           []string
	discoverSent  bool
	discoverDelay time.Duration
}

func New(ctx context.Context, addr string) (m *mqtt, err error) {
	m = &mqtt{
		addr:          addr,
		discoverDelay: 5 * time.Second,
	}
	client, err := libmqtt.NewClient(
		libmqtt.WithServer(addr),
		libmqtt.WithKeepalive(10, 1.2),
		libmqtt.WithLog(libmqtt.Info),
		libmqtt.WithRouter(newRouter(m)),
		libmqtt.WithDialTimeout(5),
	)
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
	m.initCh = make(chan error)
	m.client = client
	m.client.Connect(m.onConnect)
	err, _ = <-m.initCh

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
	if m.initCh != nil {
		if err != nil {
			m.initCh <- err
		} else {
			close(m.initCh)
		}
		m.initCh = nil
	} else if code == libmqtt.CodeSuccess {
		m.client.Publish(&libmqtt.PublishPacket{TopicName: "devnull"})
	}

	if err != nil {
		log.Printf("MQTT Connect Error: %s (%x) %v", server, code, err)
		return
	}

	if m.deviceState != nil {
		m.sendDiscover()
	}
}

func (m *mqtt) sendDiscover() {
	m.discoverSent = false
	m.client.Subscribe(&libmqtt.Topic{Name: "announce/#"})
	time.AfterFunc(m.discoverDelay, func() {
		m.discoverSent = true
		m.client.Publish(&libmqtt.PublishPacket{
			TopicName: "discover",
			IsRetain:  true,
			Payload:   []byte("1"),
		})
	})
}
