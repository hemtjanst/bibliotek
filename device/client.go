package device

import (
	"encoding/json"
	"github.com/hemtjanst/bibliotek/feature"
)

// Client device is used by applications that are talking to the actual devices,
// For example an application that's controlling lights over z-wave
type Client interface {
	// Device contains the common functions for Client and Server
	Device

	// Feature returns the feature with client functions (OnSet(), Update()) available.
	// Fetching a feature that doesn't exist will _NOT_ return nil, but instead create
	// a *feature.Fake. To check if feature actually exists; call Feature("name").Exists()
	Feature(name string) feature.Client

	// Features returns a slice of all available features
	Features() []feature.Client
}

// NewClient will create a new client device from the device.Info.
// It spawns off a goroutine that checks for inbound discover-messages
// and returns meta-data to the announce-topic.
// The transport is responsible for closing the channel, at that point the
// goroutine will stop.
func NewClient(info *Info, transport Transport) (Client, error) {
	d := &clientDev{
		device: device{
			info:      info,
			features:  map[string]feature.Feature{},
			transport: transport,
		},
	}
	err := create(&d.device)
	if err != nil {
		return nil, err
	}
	go func() {
		ch := transport.Discover()
		for {
			_, open := <-ch
			if !open {
				return
			}
			meta, _ := json.Marshal(d.info)
			transport.PublishMeta(d.Id(), meta)
		}
	}()
	return d, nil
}

type clientDev struct {
	device
}

func (d *clientDev) Feature(name string) feature.Client {
	return d.getFeature(name)
}

func (d *clientDev) Features() []feature.Client {
	var fts []feature.Client
	for _, ft := range d.features {
		fts = append(fts, ft)
	}
	return fts
}
