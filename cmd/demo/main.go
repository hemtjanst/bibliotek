// demo creates a bunch of fake devices defined in a config file in order to demo
// the hemtjänst system. It should not be used on a production system.
//
// Pass in -demo.config=./devices.json to change the config file
//
// An example config can be found in the source directory
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"lib.hemtjan.st/client"
	"lib.hemtjan.st/device"
	"lib.hemtjan.st/transport/mqtt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	flgConfig    = flag.String("demo.config", "./devices.json", "List of devices to create")
	errNoDevices = errors.New("unable to create any devices")
)

// Config extends the device.Info and adds an Init field.
type Config struct {
	Devices []struct {
		*device.Info
		Init *map[string]string `json:"init"`
	} `json:"devices"`
}

func main() {
	// Set up flags
	mCfg := mqtt.MustFlags(flag.String, flag.Bool)
	flag.Parse()

	// Set up context and signal interrupts
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Read config
	c, err := readCfg()
	if err != nil {
		log.Fatal(err)
	}

	// Create transport
	tr, err := mqtt.New(ctx, mCfg())
	if err != nil {
		log.Fatal(err)
	}

	// Loop through config and create the devices
	wg := sync.WaitGroup{}
	for _, info := range c.Devices {
		d, err := client.NewDevice(info.Info, tr)
		if err != nil {
			log.Printf("Error creating device: %v", err)
			continue
		}

		// Loop through the features of the newly created device
		for _, ft := range d.Features() {
			wg.Add(1)
			// fix mutation
			ft := ft
			d := d
			// Set up a goroutine per feature that listens for set commands
			go func() {
				defer wg.Done()
				sch, _ := ft.OnSet()
				for {
					select {
					case nv, open := <-sch:
						if !open {
							return
						}
						// Echo back whatever value was received as
						// an update.
						ft.Update(nv)
						log.Printf("OnSet(%s(%s), %s) = %s", d.Id(), d.Name(), ft.Name(), nv)
					case <-ctx.Done():
						return
					}
				}
			}()
		}

		// Set initial values from JSON for the devices
		if info.Init != nil {
			for ft, v := range *info.Init {
				d.Feature(ft).Update(v)
			}
		}
	}

	// Wait for all goroutines to end, or exit if no devices were started
	wg.Wait()
}

func readCfg() (*Config, error) {
	c := &Config{}
	f, err := ioutil.ReadFile(*flgConfig)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(f, c); err != nil {
		return nil, err
	}
	if len(c.Devices) == 0 {
		return nil, errNoDevices
	}
	return c, nil
}
