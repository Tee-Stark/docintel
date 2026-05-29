package rest

import (
	"database/sql"
	"docintel/internal/domain"

	"github.com/go-chi/chi/v5"
)

type AppServer struct {
	db          *sql.DB
	cache       domain.Cache
	repo        domain.Repository
	authService domain.AuthService
	docService  domain.DocumentService

	Router *chi.Mux
}

func NewAppServer(db *sql.DB, cache domain.Cache, repo domain.Repository, authService domain.AuthService, docService domain.DocumentService, router *chi.Mux) *AppServer {

	return &AppServer{
		db:          db,
		cache:       cache,
		repo:        repo,
		authService: authService,
		docService:  docService,
		Router:      router,
	}
}
