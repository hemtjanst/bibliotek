package device

// Device contains the common functions for Client and Server
type Common interface {
	// Id will return the unique id of the device.
	// This is currently always the same as the topic name.
	Id() string

	// Name returns the name of the device
	Name() string

	// Manufacturer returns the manufacturer
	Manufacturer() string

	// Model returns the model name/number
	Model() string

	// SerialNumber returns the serial number
	SerialNumber() string

	// Type returns the type of the device (lightbulb, outlet, etc)
	Type() string
}

type Transport interface {
	Publish(topic string, payload []byte, retain bool)
	Subscribe(topic string) chan []byte
	Unsubscribe(topic string) bool
	Resubscribe(oldTopic, newTopic string) bool
	Discover() chan struct{}
	PublishMeta(topic string, payload []byte)
	LastWillID() string
}

type Action string

const (
	DeleteAction Action = "delete"
	UpdateAction Action = "update"
	LeaveAction  Action = "leave"
)

type State struct {
	Device *Info
	Action Action
	Topic  string
}
