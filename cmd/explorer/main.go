package main

import (
	"context"
	"github.com/hemtjanst/bibliotek/server"
	"github.com/hemtjanst/bibliotek/transport/mqtt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"flag"
)

var (
	flgMqttAddr = flag.String("mqtt.addr", "localhost:1883", "Address to MQTT")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit
		cancel()
	}()

	mq, err := mqtt.New(ctx, *flgMqttAddr)

	log.Printf("Connected!")

	if err != nil {
		log.Fatal(err)
	}

	err = server.New(mq).Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
