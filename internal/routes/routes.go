package routes

import (
	"net/http"

	"docintel/internal/transport/rest/handlers"

	"github.com/go-chi/chi/v5"
)

func UserRoutes(r chi.Router) {
	r.Get("/alive", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to my document intelligence app!"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"pong"}`))
	})
}

func DocumentRoutes(r chi.Router, docHandler *handlers.DocumentHandler) {
	r.Post("/documents/upload", docHandler.UploadDocument)
}
