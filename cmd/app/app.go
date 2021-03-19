package main

import (
	"context"
	"log"
	"runtime"

	"github.com/ftamas88/kafka-test/internal/app"
)

// main Initialises the App
func main() {
	a, err := app.New()
	if err != nil {
		log.Panicf("unable to start application: %s", err.Error())
	}
	defer a.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	go app.SignalShutdownHandler(cancelFunc)

	err = a.Run(ctx)
	if err != nil {
		log.Printf("Fatal error: %s", err)
	}

	runtime.Goexit()
}
