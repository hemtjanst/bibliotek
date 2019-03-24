package main

import (
	"context"
	"flag"
	"github.com/hemtjanst/bibliotek/device"
	"github.com/hemtjanst/bibliotek/server"
	"github.com/hemtjanst/bibliotek/transport/mqtt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mCfg := mqtt.MustFlags(flag.String, flag.Bool)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit
		cancel()
	}()

	mq, err := mqtt.New(ctx, mCfg())

	log.Printf("Connected!")

	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(mq)

	srv.SetHandler(&handler{})

	err = srv.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

type handler struct{}

func (h *handler) AddedDevice(dev server.Device) {
	log.Printf("Device Created: %+v", dev)

}

func (h *handler) UpdatedDevice(dev server.Device, updates []*device.InfoUpdate) {
	for _, upd := range updates {
		log.Printf("[%s] %s changed \"%s\" -> \"%s\" (%+v)",
			dev.Id(),
			upd.Field,
			upd.Old,
			upd.New,
			upd.FeatureInfo,
		)
	}

}

func (h *handler) RemovedDevice(dev server.Device) {
	log.Printf("Device Removed: %+v", dev)
}
