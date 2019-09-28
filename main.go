package main

import (
	"os"
	"os/signal"

	"cafebot/chat"
	"cafebot/config"
	"cafebot/order"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		fatalf("loading config: %v\n", err)
	}

	setupLog(os.Stderr, conf.LogLevel)

	orderService := order.NewInMemoryService(conf)

	chat, err := chat.NewService(
		conf,
		orderService,
	)
	if err != nil {
		fatalf("creating chat service: %s", err)
	}

	err = chat.Run()
	if err != nil {
		fatalf("running chat service: %s", err)
	}

	// Shutdown on SIGINT (CTRL-C).
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	chat.Stop()

	log.Infof("shutting down...")
}
