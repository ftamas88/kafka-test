package routing

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ftamas88/kafka-test/internal/controller"
	"github.com/gorilla/mux"
)

type RouterConfig struct {
	WellKnown *controller.WellKnownController
}

type Router struct {
	httpPort int
	handler  http.Handler
}

func NewRouter(httpPort int, conf *RouterConfig) *Router {
	r := mux.NewRouter()

	r.HandleFunc("/", conf.WellKnown.Root).Methods(http.MethodGet)

	status := r.PathPrefix("/status").Subrouter()
	status.HandleFunc("/health", conf.WellKnown.Health).Methods(http.MethodGet)
	status.HandleFunc("/version", conf.WellKnown.Version).Methods(http.MethodGet)

	return &Router{
		httpPort: httpPort,
		handler:  r,
	}
}

func (r *Router) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", r.httpPort))
	if err != nil {
		return fmt.Errorf("router net.Listen: %w", err)
	}

	server := http.Server{
		Handler: r.handler,
	}

	go func() {
		//<-ctx.Done()
		//err := server.Shutdown(ctx)
		//if err != nil {
		//	log.Fatalf("error shutting down router: %s", err)
		//}
	}()

	log.Printf("HTTP server running on port %d", r.httpPort)

	err = server.Serve(listener)
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
