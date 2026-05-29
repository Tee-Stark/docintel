package main

import (
	"docintel/internal/adapters/postgres"
	"docintel/internal/adapters/redis"
	"docintel/internal/app"
	"docintel/internal/transport/rest"
	"docintel/pkg/config"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	db, err := config.NewDBConfig().ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	repository := postgres.NewRepository(db)

	redisClient, err := config.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	cache := redis.NewCache(redisClient)

	authService := app.NewAuthService(repository, cache)
	docService := app.NewDocumentService(repository)

	r := chi.NewRouter()
	app := rest.NewAppServer(db, cache, repository, authService, docService, r)
	StartServer(app)
}

func StartServer(app *rest.AppServer) {
	rest.SetUpRoutes(app)

	log.Println("Server running on port 8080")

	if err := http.ListenAndServe(":8080", app.Router); err != nil {
		log.Fatal(err)
	}
}
