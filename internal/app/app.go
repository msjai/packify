package app

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"packify/internal/config"
	"packify/internal/httpserver"
	"packify/internal/logger"
)

func Run() {
	cfg := config.MustLoad()

	logger.New(cfg.Env)

	slog.Info("starting packify server", slog.String("env", cfg.Env))

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	})

	srv := httpserver.New(cfg.HTTPAddress, router)
	srv.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("signal received", "signal", s.String())
	case err := <-srv.Notify():
		slog.Error("server error", "error", err)
	}

	if err := srv.Shutdown(); err != nil {
		slog.Error("shutdown error", "error", err)
	}

	slog.Info("server stopped")
}
