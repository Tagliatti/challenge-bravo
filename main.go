package main

import (
	"github.com/Tagliatti/challenge-bravo/database"
	"github.com/Tagliatti/challenge-bravo/handler/health"
	"log"
	"net/http"
)

func main() {
	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Database migrating...")
	database.Migrate(db)

	healthHandler := health.NewHealthHandler()

	http.HandleFunc("GET /health", healthHandler.Handle)

	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
