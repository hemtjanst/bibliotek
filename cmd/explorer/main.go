package main

import (
	"context"
	"flag"
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

	err = server.New(mq).Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
