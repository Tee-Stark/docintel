package main

import (
	"docintel/internal/routes"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	r := chi.NewRouter()

	routes.UserRoutes(r)

	log.Println("Server running on port 8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
