package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/url"
	"slices"
	"sync"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"

	"lib.hemtjan.st/v2/component"
	"lib.hemtjan.st/v2/device"
)

type Server struct {
	Devices []*device.Device

	subscribe []string

	logger     *slog.Logger
	reqTimeout time.Duration

	pahoConfig autopaho.ClientConfig
	pahoMgr    *autopaho.ConnectionManager
	pahoRouter *paho.StandardRouter

	sync.RWMutex
}

func (s *Server) Subscribe(ctx context.Context, topic string, handler paho.MessageHandler) error {
	s.Lock()
	existing := slices.Contains(s.subscribe, topic)
	if !existing {
		s.subscribe = append(s.subscribe, topic)
	}
	s.Unlock()

	s.pahoRouter.RegisterHandler(topic, handler)
	if !existing && s.pahoMgr != nil {
		rctx, cancel := timeout(ctx, s.reqTimeout)
		defer cancel()
		_, err := s.pahoMgr.Subscribe(rctx, &paho.Subscribe{Subscriptions: []paho.SubscribeOptions{
			{Topic: topic},
		}})
		return err
	}
	return nil
}

func (s *Server) Publish(ctx context.Context, topic string, qos uint8, msg []byte) error {
	rctx, cancel := timeout(ctx, s.reqTimeout)
	defer cancel()
	_, err := s.pahoMgr.Publish(rctx,
		&paho.Publish{
			QoS:     byte(qos),
			Topic:   topic,
			Payload: msg,
		},
	)

	return err
}

func (s *Server) AddDevice(ctx context.Context, device *device.Device) error {
	s.Lock()
	s.Devices = append(s.Devices, device)
	cm := s.pahoMgr
	s.Unlock()

	if device.Components == nil {
		return nil
	}

	for _, cmp := range device.Components {
		if cmpBase, ok := cmp.(component.BaseComponent); ok {
			cmpRef := cmpBase.GetBaseReference()
			if cmpRef.AvailabilityTopic == "" && len(cmpRef.Availability) == 0 {
				cmpRef.AvailabilityTopic = s.WillTopic()
			}
		}

		if cmpUpdatable, ok := cmp.(component.Updatable); ok {
			for _, c := range cmpUpdatable.UpdateChannels() {
				go func(c component.UpdateChannel) {
					for {
						msg, open := <-c.Channel
						if !open {
							return
						}
						_ = s.Publish(ctx, c.Topic, 1, []byte(msg))
					}
				}(c)
			}
		}
		if cmpCommandable, ok := cmp.(component.Commandable); ok {
			for _, c := range cmpCommandable.CommandChannels() {
				c := c
				_ = s.Subscribe(ctx, c.Topic, func(publish *paho.Publish) {
					c.Channel <- string(publish.Payload)
				})
			}
		}
	}

	if cm != nil {
		if err := s.publishDevice(ctx, s.pahoMgr, device); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) publishDevice(ctx context.Context, cm *autopaho.ConnectionManager, device *device.Device) error {
	buf, err := json.Marshal(device)
	if err != nil {
		return err
	}

	rctx, cancel := timeout(ctx, s.reqTimeout)
	defer cancel()
	_, err = cm.Publish(
		rctx,
		&paho.Publish{
			QoS:     byte(1),
			Topic:   device.DiscoveryTopic(),
			Payload: buf,
		},
	)

	return err
}

func (s *Server) Start(ctx context.Context) error {
	c, err := autopaho.NewConnection(ctx, s.pahoConfig)
	if err != nil {
		return err
	}

	if err = c.AwaitConnection(ctx); err != nil {
		return err
	}

	s.pahoMgr = c

	return s.Subscribe(ctx, "homeassistant/status", func(publish *paho.Publish) {
		if string(publish.Payload) == "online" {
			if len(s.Devices) == 0 {
				return
			}
			for _, dev := range s.Devices {
				if err = s.publishDevice(ctx, s.pahoMgr, dev); err != nil {
					s.logger.Error("unable to publish device", slog.String("error", err.Error()))
				}
			}
		}
	})
}

func (s *Server) WillTopic() string {
	return "homeassistant/client/" + s.pahoConfig.ClientID + "/status"
}

func New(ctx context.Context, log *slog.Logger, u string, clientID string) (*Server, error) {
	srv, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	r := paho.NewStandardRouterWithDefault(func(p *paho.Publish) {})

	var s *Server
	s = &Server{
		reqTimeout: 5 * time.Second,
		logger:     log,
		pahoRouter: r,
		pahoConfig: autopaho.ClientConfig{
			Errors:                        slog.NewLogLogger(log.Handler(), slog.LevelError),
			PahoErrors:                    slog.NewLogLogger(log.Handler(), slog.LevelError),
			ServerUrls:                    []*url.URL{srv},
			KeepAlive:                     20,
			CleanStartOnInitialConnection: false,
			SessionExpiryInterval:         60,
			WillMessage: &paho.WillMessage{
				QoS:     2,
				Topic:   "homeassistant/client/" + clientID + "/status",
				Payload: []byte("offline"),
			},
			OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
				s.Lock()
				devs := s.Devices
				subscr := s.subscribe
				s.Unlock()

				if len(subscr) > 0 {
					var subs []paho.SubscribeOptions
					for _, topic := range subscr {
						subs = append(subs, paho.SubscribeOptions{Topic: topic})
					}
					rctx, cancel := timeout(ctx, s.reqTimeout)
					defer cancel()
					_, err := cm.Subscribe(rctx, &paho.Subscribe{Subscriptions: subs})
					if err != nil {
						s.logger.Error("unable to subscribe", slog.String("error", err.Error()))
					}
				}

				for _, dev := range devs {
					if err = s.publishDevice(ctx, cm, dev); err != nil {
						s.logger.Error("unable to publish device", slog.String("error", err.Error()))
					}
				}

				rctx, cancel := timeout(ctx, s.reqTimeout)
				defer cancel()
				_, err = cm.Publish(rctx, &paho.Publish{
					QoS:     2,
					Topic:   "homeassistant/client/" + clientID + "/status",
					Payload: []byte("online"),
				})

				if err != nil {
					s.logger.Error("unable to publish status", slog.String("error", err.Error()))
				}

				s.logger.Info("MQTT connected")
			},
			OnConnectError: func(err error) { s.logger.Error("unable to connect", slog.String("error", err.Error())) },
			ClientConfig: paho.ClientConfig{
				ClientID:      clientID,
				Router:        r,
				OnClientError: func(err error) { s.logger.Error("client issue", slog.String("error", err.Error())) },
				OnServerDisconnect: func(d *paho.Disconnect) {
					if d.Properties != nil {
						s.logger.Info("server requested disconnect", slog.String("reason", d.Properties.ReasonString))
					} else {
						s.logger.Info("server requested disconnect", slog.Int("code", int(d.ReasonCode)))
					}
				},
			},
		},
	}

	return s, nil
}
