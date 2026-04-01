package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// BuildRoutes registers API routes and their handlers.
func BuildRoutes(router *chi.Mux,
	calculateService Calculator,
	getService Getter,
	submitService Submitter,
) {
	packCalculateHandler := NewCalculateHandler(calculateService)
	packGetHandler := NewGetHandler(getService)
	packSubmitHandler := NewSubmitHandler(submitService)

	pingHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	}
	router.Get("/ping", pingHandler)
	router.Head("/ping", pingHandler)

	// Group allows applying middleware to these routes in the future.
	router.Group(func(router chi.Router) {
		router.Get("/api/packs", packGetHandler.Handle)
		router.Put("/api/packs", packSubmitHandler.Handle)
		router.Post("/api/calculate", packCalculateHandler.Handle)
	})

}
