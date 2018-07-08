package mqtt

import "github.com/goiiot/libmqtt"

type mqttClient interface {
	// Connect to all specified server with client options
	Connect(handler libmqtt.ConnHandler)

	// Publish a message for the topic
	Publish(packets ...*libmqtt.PublishPacket)

	// Subscribe topic(s)
	Subscribe(topics ...*libmqtt.Topic)

	// UnSubscribe topic(s)
	UnSubscribe(topics ...string)

	// Destroy all client connection
	Destroy(force bool)
}
