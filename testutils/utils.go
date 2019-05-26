package testutils // import "lib.hemtjan.st/testutils"

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/goiiot/libmqtt"
	"lib.hemtjan.st/client"
	"lib.hemtjan.st/device"
)

var errNoDevices = errors.New("unable to create any devices")

type dummyMqtt struct {
	initCh chan error
}

// onConnect handler
func (d *dummyMqtt) onConnect(server string, code byte, err error) {
	if code != libmqtt.CodeSuccess && err == nil {
		err = fmt.Errorf("error code %d", int(code))
	}

	if d.initCh != nil {
		if err != nil {
			d.initCh <- err
		} else {
			close(d.initCh)
		}
	}
}

// MQTTAddress is a helper function that will try to connect
// to an MQTT broker on MQTT_ADDRESS and if successful return
// the address.
// It should only be used in Bibliotek's own tests or other
// Hemtjanst component tests
func MQTTAddress(t *testing.T) string {
	t.Helper()

	mqttHostPort := os.Getenv("MQTT_ADDRESS")
	if mqttHostPort == "" {
		mqttHostPort = "localhost:1883"
	}

	client, err := libmqtt.NewClient(
		libmqtt.WithServer(mqttHostPort),
		libmqtt.WithDialTimeout(1),
	)

	if err != nil {
		t.Fatalf(err.Error())
	}

	d := &dummyMqtt{
		initCh: make(chan error),
	}

	client.Connect(d.onConnect)
	err, _ = <-d.initCh
	d.initCh = nil

	if err != nil {
		client.Destroy(true)
		t.Fatalf(err.Error())
	}

	defer client.Destroy(false)

	return mqttHostPort
}

type config struct {
	Devices []struct {
		*device.Info
		Init *map[string]string `json:"init"`
	} `json:"devices"`
}

// DevicesFromJSON creates fake devices on the broker based
// on the devices defined in a JSON file
// When you run the returned cleanup function it will delete the devices.
// Note that unless you then also close the transport it might still
// announce the device(s) when it receives a discover.
func DevicesFromJSON(path string, m device.Transport) (func(), error) {
	c := &config{}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return func() {}, err
	}
	if err = json.Unmarshal(f, c); err != nil {
		return func() {}, err
	}
	if len(c.Devices) == 0 {
		return func() {}, errNoDevices
	}
	// Loop through config and create the devices
	for _, info := range c.Devices {
		d, err := client.NewDevice(info.Info, m)
		if err != nil {
			log.Printf("Error creating device: %v", err)
			continue
		}
		if info.Init != nil {
			for ft, v := range *info.Init {
				d.Feature(ft).Update(v)
			}
		}
	}

	cleanup := func() {
		for _, info := range c.Devices {
			client.DeleteDevice(info.Info, m)
		}
	}
	return cleanup, nil
}
