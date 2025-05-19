package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lib.hemtjan.st/component"
	"lib.hemtjan.st/device"
	"lib.hemtjan.st/server"
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

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	)

	srv := check1(server.New(ctx, logger, os.Getenv("MQTT_ADDRESS"), "bibliotek-demo-sensor"))
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
	check(srv.AddDevice(ctx, dev))

	<-ctx.Done()
	stop()
}
