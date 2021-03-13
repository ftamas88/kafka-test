package app

import (
	"context"
	"log"

	"golang.org/x/sync/errgroup"
)

type App struct {
	version    string
	commitHash string
	router     Component
	consumer   Component
	producer	Component
	cloudProducer Component
	watcher 	Component
	closeFuncs []Close
}

var version string
var commitHash string

// Version returns the git tag used to build the application
func Version() string {
	return version
}

// CommitHash returns the git commit hash used to build the application
func CommitHash() string {
	return commitHash
}

// Run spawns go routines for the each of its components, starts shutting down all components if one fails
// or the context is closed.
//
// Returns the first error encountered or nil once all components shutdown gracefully.
func (app App) Run(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return app.router.Run(egCtx)
	})

	eg.Go(func() error {
		return app.producer.Run(egCtx)
	})

	eg.Go(func() error {
		return app.cloudProducer.Run(egCtx)
	})

	eg.Go(func() error {
		return app.watcher.Run(egCtx)
	})

	eg.Go(func() error {
		return app.consumer.Run(egCtx)
	})

	log.Printf("Application is running version: %s", Version())

	return eg.Wait()
}
