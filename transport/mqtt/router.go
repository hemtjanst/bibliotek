package mqtt

import (
	"github.com/goiiot/libmqtt"
	"strings"
)

type mqttRouter struct {
	handler MessageHandler
}

func newRouter(handler MessageHandler) libmqtt.TopicRouter {
	return &mqttRouter{
		handler: handler,
	}
}

// Name is the name of router
func (r *mqttRouter) Name() string {
	return "Hemtjänst MQTT Router"
}

// Handle defines how to register topic with handler
func (r *mqttRouter) Handle(topic string, h libmqtt.TopicHandler) {

}

// Dispatch defines the action to dispatch published packet
func (r *mqttRouter) Dispatch(p *libmqtt.PublishPacket) {
	if p.TopicName == leaveTopic {
		r.handler.OnLeave(p)
		return
	}
	if p.TopicName == discoverTopic {
		r.handler.OnDiscover(p)
		return
	}
	if p.TopicName == announceTopic || strings.HasPrefix(p.TopicName, announceTopic+"/") {
		r.handler.OnAnnounce(p)
		return
	}
	r.handler.OnFeature(p)
}
