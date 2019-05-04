package mqtt

import "lib.hemtjan.st/device"

type MQTT interface {
	TopicName(t EventType) string
	DeviceState() chan *device.State
	PublishMeta(topic string, payload []byte)
	Publish(topic string, payload []byte, retain bool)
	SubscribeRaw(topic string) chan *Packet
	Unsubscribe(topic string) bool
	Resubscribe(oldTopic, newTopic string) bool
	Subscribe(topic string) chan []byte
	Discover() chan struct{}
	LastWillID() string
}
