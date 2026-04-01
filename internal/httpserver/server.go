package httpserver

import (
	"context"
	"net/http"
	"time"
)

// Server wraps http.Server with graceful shutdown support.
// The notify channel reports server errors without blocking the main goroutine.
type Server struct {
	server *http.Server
	notify chan error
}

// New creates a new HTTP server with the given address and handler.
// Timeouts protect against slow clients (slowloris attack).
func New(addr string, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
		},
		notify: make(chan error, 1),
	}
}

// Start runs the server in a separate goroutine.
// Any error (e.g. port already in use) is sent to the notify channel.
func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify returns a read-only channel for server errors.
// Used in select to react to server failures.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown gracefully stops the server.
// Allows 5 seconds for in-flight requests to complete, then forces shutdown.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
