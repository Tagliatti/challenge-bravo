package main

import (
	"github.com/Tagliatti/challenge-bravo/database"
	"github.com/Tagliatti/challenge-bravo/handler/conversion"
	"github.com/Tagliatti/challenge-bravo/handler/currency"
	"github.com/Tagliatti/challenge-bravo/handler/health"
	"github.com/Tagliatti/challenge-bravo/repository"
	"github.com/Tagliatti/challenge-bravo/service"
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

	currencyRepository := repository.NewCurrencyRepository(db)
	currencyApi := service.NewCurrencyApi(http.DefaultClient)

	healthHandler := health.NewHealthHandler()
	createHandler := currency.NewCreateHandler(currencyRepository, currencyApi).Handler
	updateHandler := currency.NewUpdateHandler(currencyRepository, currencyApi).Handler
	deleteHandler := currency.NewDeleteHandler(currencyRepository).Handler
	conversionHandler := conversion.NewConversionHandler(currencyRepository, currencyApi).Handler

	http.HandleFunc("GET /health", healthHandler.Handle)
	http.HandleFunc("POST /", createHandler)
	http.HandleFunc("PUT /{symbol}", updateHandler)
	http.HandleFunc("DELETE /{symbol}", deleteHandler)
	http.HandleFunc("GET /", conversionHandler)

	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
