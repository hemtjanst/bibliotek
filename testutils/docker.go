package testutils // import "lib.hemtjan.st/testutils"

import (
	"fmt"
	"log"

	"github.com/goiiot/libmqtt"
	"github.com/ory/dockertest"
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

// mqttBroker starts an MQTT broker using Docker
// It returns a host:port string as well as a cleanup function that you can defer
func MQTTBroker() (string, func()) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("ansi/mosquitto", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	host := resource.GetHostPort("1883/tcp")
	if host == "" {
		log.Fatalf("Empty host")
	}

	if err := pool.Retry(func() error {
		client, err := libmqtt.NewClient(
			libmqtt.WithServer(host),
			libmqtt.WithDialTimeout(1),
		)
		if err != nil {
			log.Fatal(err)
		}
		d := &dummyMqtt{
			initCh: make(chan error),
		}
		client.Connect(d.onConnect)
		err, _ = <-d.initCh
		d.initCh = nil

		if err != nil {
			client.Destroy(true)
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to MQTT: %s", err)
	}

	return host, func() {
		if err := pool.Purge(resource); err != nil {
			log.Printf("Could not purge resource: %s", err)
		}
	}
}
