package testutils // import "lib.hemtjan.st/testutils"

import (
	"fmt"
	"os"
	"testing"

	"github.com/goiiot/libmqtt"
)

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
