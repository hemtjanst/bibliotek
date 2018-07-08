package client

type Client struct {
}

type Transport interface {
	// HandleDeviceState should get a callback when a discover query is made
	Discover() chan bool

	/*
		// HandleFeatureState is called with the id of the device and a feature name,
		// whenever that feature updates - the callback is called with the new value
		FeatureState(device, feature string) chan string

		// SetState should update the value of a device feature
		SetState(id, feature, newState string) error
	*/
}

func New() *Client {
	return &Client{}
}
