package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Close interface {
	Close() error
}

// Close loops through all cancel/close functions and calls them on shutdown
func (app App) Close() {
	for _, f := range app.closeFuncs {
		_ = f.Close()
	}
}

// SignalShutdownHandler listens for a SIGTERM signal
// and gracefully cancels the main application context
// once this is completed exits the app
func SignalShutdownHandler(cancelFunction context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// Invoke the cancel function
	cancelFunction()

	// Safety deadline for exiting
	<-time.After(10 * time.Second)
	os.Exit(1)
}
