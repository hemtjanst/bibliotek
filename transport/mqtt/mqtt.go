package mqtt

import (
	"context"
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
	errCh         chan error
	stopCh        []chan struct{}
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
	ctx           context.Context
	sync.RWMutex
}

func New(ctx context.Context, c *Config) (MQTT, error) {
	if c == nil {
		return nil, ErrNoConfig
	}
	if err := c.check(); err != nil {
		return nil, err
	}
	m := &mqtt{
		discoverDelay: c.DiscoverDelay,
		willID:        c.ClientID,
		announceTopic: c.AnnounceTopic,
		discoverTopic: c.DiscoverTopic,
		leaveTopic:    c.LeaveTopic,
		willMap:       map[string][]string{},
		ctx:           ctx,
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

	opts := []libmqtt.Option{
		libmqtt.WithRouter(newRouter(m)),
	}
	opts = append(opts, c.opts()...)

	client, err := libmqtt.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	m.client = client

	if ctx != nil && ctx.Done() != nil {
		go func() {
			<-ctx.Done()
			m.destroy()
		}()
	}

	return m, nil
}

func (m *mqtt) isCancelled() bool {
	if m.ctx != nil {
		select {
		case <-m.ctx.Done():
			return true
		default:
		}
	}
	return false
}

type Err string

func (e Err) Error() string { return string(e) }

const (
	ErrIsAlreadyRunning Err = "already running"
	ErrIsCancelled      Err = "cancelled"
)

func (m *mqtt) Start() (bool, error) {
	m.Lock()
	if m.errCh != nil {
		m.Unlock()
		return false, ErrIsAlreadyRunning
	}
	if m.isCancelled() {
		m.Unlock()
		return false, ErrIsCancelled
	}
	errCh := make(chan error)
	m.errCh = errCh
	if m.sub == nil {
		m.sub = map[string][]chan []byte{}
	}
	if m.subRaw == nil {
		m.subRaw = map[string][]chan *Packet{}
	}
	cl := m.client
	m.Unlock()

	defer func() {
		m.Lock()
		stopCh := m.stopCh
		m.stopCh = nil
		m.errCh = nil
		m.Unlock()
		for _, ch := range stopCh {
			close(ch)
		}
	}()
	if cl == nil {
		return false, ErrIsCancelled
	}

	cl.Connect(m.onConnect)
	err := <-errCh
	return !m.isCancelled(), err
}

func (m *mqtt) destroy() {
	stopCh := make(chan struct{})
	m.Lock()
	dsub := m.discoverSub
	subs := m.sub
	subRaw := m.subRaw
	stateCh := m.deviceState
	cl := m.client
	if m.errCh != nil {
		m.stopCh = append(m.stopCh, stopCh)
	} else {
		close(stopCh)
	}
	if m.errCh != nil {
		close(m.errCh)
		m.errCh = nil
	}
	m.discoverSub = []chan struct{}{}
	m.sub = map[string][]chan []byte{}
	m.subRaw = map[string][]chan *Packet{}
	m.deviceState = nil
	m.client = nil
	m.Unlock()

	if cl != nil {
		cl.Destroy(true)
	}

	if stateCh != nil {
		close(stateCh)
	}

	for _, ch := range dsub {
		close(ch)
	}

	if subRaw != nil {
		for _, v := range subRaw {
			for _, vv := range v {
				func() {
					defer func() {
						_ = recover()
					}()
					close(vv)
				}()
			}
		}
	}

	if subs != nil {
		for _, chans := range subs {
			for _, ch := range chans {
				func() {
					defer func() {
						_ = recover()
					}()
					close(ch)
				}()
			}
		}
	}
}

func (m *mqtt) onConnect(server string, _code byte, err error) {
	m.Lock()
	defer m.Unlock()
	code := Code(_code)

	if code != libmqtt.CodeSuccess && err == nil {
		err = code
	}

	if err != nil {
		if m.errCh != nil {
			m.errCh <- err
			return
		}
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
		if m.client != nil {
			m.discoverSent = true
			m.client.Publish(&libmqtt.PublishPacket{
				TopicName: m.discoverTopic,
				IsRetain:  true,
				Payload:   []byte("1"),
			})
		}
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
