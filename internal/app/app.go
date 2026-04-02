package app

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/msjai/packify/internal/config"
	"github.com/msjai/packify/internal/handler"
	"github.com/msjai/packify/internal/httpserver"
	"github.com/msjai/packify/internal/logger"
	"github.com/msjai/packify/internal/service"
	"github.com/msjai/packify/internal/storage/sqlite"
	"github.com/msjai/packify/web"
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

	handler.BuildRoutes(router, calculateService, getService, submitService, cfg.MaxOrder)

	// Serve embedded static files (index.html, app.js, etc.).
	router.Handle("/*", http.FileServerFS(web.FS))

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
