package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	notify chan error
}

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

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
