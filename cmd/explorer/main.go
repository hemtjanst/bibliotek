package main

import (
	"context"
	"flag"
	"github.com/hemtjanst/bibliotek/client"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hemtjanst/bibliotek/server"
	"github.com/hemtjanst/bibliotek/transport/mqtt"
)

func main() {
	mCfg := mqtt.MustFlags(flag.String, flag.Bool)
	flgDelete := flag.String("delete", "", "Delete device (by topic)")
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

	ch := make(chan server.Update, 1)
	srv.SetUpdateChannel(ch)
	go func() {
		for {
			upd, open := <-ch
			if !open {
				return
			}
			dev := upd.Device
			switch upd.Type {
			case server.AddedDevice:
				if *flgDelete != "" && dev.Id() == *flgDelete {
					_ = client.DeleteDevice(dev.Info(), mq)
				}
				log.Printf("Device Created: %s (\"%s\" with serial %s model %s made by %s)",
					dev.Id(), dev.Name(), dev.SerialNumber(), dev.Model(), dev.Manufacturer())
				for _, ft := range dev.Features() {
					log.Printf("                * %s (min/max/step: %d/%d/%d", ft.Name(), ft.Min(), ft.Max(), ft.Step())
					go func(ft server.Feature) {
						ch, _ := ft.OnUpdate()
						for {
							d, open := <-ch
							if !open {
								log.Printf("Device: %s feature: %s disappeared", dev.Id(), ft.Name())
								return
							}
							log.Printf("Device: %s feature: %s is now %s", dev.Id(), ft.Name(), d)
						}
					}(ft)
				}
			case server.UpdatedDevice:
				for _, upd := range upd.Changes {
					log.Printf("[%s] %s changed \"%s\" -> \"%s\" (%+v)",
						dev.Id(),
						upd.Field,
						upd.Old,
						upd.New,
						upd.FeatureInfo,
					)
				}
			case server.RemovedDevice:
				log.Printf("Device Removed: %+v", dev)
			}

		}
	}()

	err = srv.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
