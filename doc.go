// Bibliotek is a library for the Hemtjänst ecosystem.
//
// It comes with a few different packages, most importantly
// client, server and transport.
//
// Client
//
// The client package lets you create new devices and register
// them onto MQTT, so that other things can become aware of it.
// For example, if you wanted to take devices from an IoT gateway
// like the IKEA Tradfri gateway and make them available as
// Hemtjänst devices, this is where you'd start.
//
// Server
//
// The server package lets you observe all existing devices,
// subscribe to state changes as well as change the state of
// existing devices. This is useful if you want to create a
// bridge to somewhere else. A HomeKit bridge could be implemented
// with it, or exposing device information through another
// medium, like providing an HTTP API.
//
// Transport
//
// The transport package contains some MQTT related utilities.
// You'll likely never need them, aside from transport/mqtt.Flags()
// so you don't have to define all the different flags for setting
// up a connection to an MQTT broker for ya CLI utility yourself.

package bibliotek
