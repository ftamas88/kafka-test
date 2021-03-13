package app

import "context"

// Component is part of the app which runs in a goroutine
type Component interface {
	// Run should cleanup the component if the context is done, and should wait for any current
	// tasks to complete before returning an error if there was one.
	Run(ctx context.Context) error
}
