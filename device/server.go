package device

import "github.com/hemtjanst/bibliotek/feature"

// Server device is used applications that monitor or send commands to the
// different client devices.
// For examples an application that integrates with a personal assistant
type Server interface {
	// Device contains the common functions for Client and Server
	Device

	// Feature returns the feature with server functions (Set(), OnUpdate()) available.
	// Fetching a feature that doesn't exist will _NOT_ return nil, but instead create
	// a *feature.Fake. To check if feature actually exists; call Feature("name").Exists()
	Feature(name string) feature.Server

	// Exists returns true if the device exists
	Exists() bool

	// Features returns a slice of available features
	Features() []feature.Server
}

// NewServer should normally only be called with data from announcements.
func NewServer(info *Info, transport Transport) (Server, error) {
	d := &serverDev{
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
	return d, nil
}

type serverDev struct {
	device
}

func (d *serverDev) Feature(name string) feature.Server {
	return d.getFeature(name)
}

func (d *serverDev) Features() []feature.Server {
	var fts []feature.Server
	for _, ft := range d.features {
		fts = append(fts, ft)
	}
	return fts
}
