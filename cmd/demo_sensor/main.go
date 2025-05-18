package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lib.hemtjan.st/v2/component"
	"lib.hemtjan.st/v2/device"
	"lib.hemtjan.st/v2/server"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func check1[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := check1(server.New(os.Getenv("MQTT_ADDRESS"), "bibliotek-demo-sensor"))
	check(srv.Start(ctx))

	sensor := component.NewTempSensor("Temperature", "bibliotek_demo_temp")
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
				currentTemp := rand.Intn(10) + 15 // Random between 15-25
				sensor.StateCh <- fmt.Sprintf("%d", currentTemp)
			}
		}
	}()

	dev := &device.Device{
		Info: device.Info{
			ID:           "bibliotek_demo",
			Name:         "Bibliotek Demo",
			Manufacturer: "hemtjanst",
			Model:        "bibliotek",
			ModelID:      "bibliotek",
		},
		Origin: device.Origin{
			Name: "bibliotek",
		},
	}
	check(dev.SetComponent("temp", sensor))
	check(srv.AddDevice(dev))

	<-ctx.Done()
	stop()
}
