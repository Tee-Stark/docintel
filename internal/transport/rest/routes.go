package rest

import (
	"docintel/internal/transport/rest/handlers"
	"docintel/internal/transport/rest/middleware"

	"github.com/go-chi/chi/v5"
)

func SetUpRoutes(app *AppServer) {
	app.Router.Use(middleware.CORS)

	middleware := middleware.NewMiddleWare(app.authService)

	// API routes
	app.Router.Route("/api/v1", func(r chi.Router) {
		// Authentication routes (public)
		userHandler := handlers.NewUserHandler(app.authService)
		r.Group(func(r chi.Router) {
			r.Post("/register", userHandler.HandleRegister)
			r.Post("/login", userHandler.HandleLogin)
		})

		// Document routes (protected)
		docHandler := handlers.NewDocumentHandler(app.docService)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authorize)
			r.Post("/documents/upload", docHandler.UploadDocument)
		})
	})
}
