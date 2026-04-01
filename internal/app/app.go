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
	"packify/internal/handler"
	"packify/internal/httpserver"
	"packify/internal/logger"
	"packify/internal/service"
	"packify/internal/storage/sqlite"
)

func Run() {
	cfg := config.MustLoad()

	logger.New(cfg.Env)
	slog.Info("starting packify server", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath, cfg.DefaultPacks)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	calculateService := service.NewCalculateService(storage)
	getService := service.NewGetService(storage)
	submitService := service.NewSubmitService(storage)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	handler.BuildRoutes(router, calculateService, getService, submitService)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	})

	srv := httpserver.New(cfg.HTTPAddress, router)
	srv.Start()
	slog.Info("server started", slog.String("address", cfg.HTTPAddress))

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
