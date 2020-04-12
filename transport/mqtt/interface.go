package mqtt

import (
	"lib.hemtjan.st/device"
)

type MQTT interface {
	// Start connects to mqtt and block until disconnected.
	// "ok" is true if the client is still valid and should be reused by calling Start() again
	//
	// Example of running with reconnect:
	//   for {
	//     ok, err := client.Start(ctx)
	//     if !ok {
	//       break
	//     }
	//     log.Printf("Error %v - retrying in 5 seconds", err)
	//     time.Sleep(5 * time.Second)
	//   }
	Start() (ok bool, err error)
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
