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
	if p.TopicName == r.handler.TopicName(TypeLeave) {
		r.handler.OnLeave(p)
		return
	}
	if p.TopicName == r.handler.TopicName(TypeDiscover) {
		r.handler.OnDiscover(p)
		return
	}
	anTopic := r.handler.TopicName(TypeAnnounce)
	if p.TopicName == anTopic || strings.HasPrefix(p.TopicName, anTopic+"/") {
		r.handler.OnAnnounce(p)
		return
	}
	r.handler.OnFeature(p)
}
