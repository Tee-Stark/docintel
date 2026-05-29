package main

import (
	"docintel/internal/routes"
	"docintel/internal/transport/http/handlers"
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

	docHandler := handlers.NewDocumentHandler(db, uploadDir)

	r := chi.NewRouter()

	routes.UserRoutes(r)
	routes.DocumentRoutes(r, docHandler)

	log.Println("Server running on port 8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
