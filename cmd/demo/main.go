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
	"fmt"
	"io/ioutil"
	"lib.hemtjan.st/client"
	"lib.hemtjan.st/device"
	"lib.hemtjan.st/hass"
	"lib.hemtjan.st/transport/mqtt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
	ctx, cancel := context.WithCancel(context.Background())
	go waitSig(cancel)

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
	go func() {
		for {
			ok, err := tr.Start()
			if err != nil {
				log.Print(err)
			}
			if !ok {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

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

		if tr.TopicName(mqtt.TypeHassPrefix) != "" {
			hd, err := hass.HTtoHA(info.Info)
			if err != nil {
				log.Printf("Error creating HASS device: %v", err)
			} else {
				for _, hdd := range hd {
					hds, _ := json.MarshalIndent(hdd, "", "  ")
					log.Print(string(hds))
					haTopic := fmt.Sprintf("%s/%s/%s", tr.TopicName(mqtt.TypeHassPrefix), hdd.Type, hdd.UniqueId)
					tr.Publish(haTopic, hds, true)
				}
			}
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

func waitSig(cancel func()) {

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	cancel()
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
