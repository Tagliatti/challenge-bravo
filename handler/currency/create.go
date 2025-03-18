package currency

import (
	"encoding/json"
	"github.com/Tagliatti/challenge-bravo/entity"
	"github.com/Tagliatti/challenge-bravo/handler"
	"github.com/Tagliatti/challenge-bravo/repository"
	"github.com/Tagliatti/challenge-bravo/service"
	"net/http"
	"strings"
)

type Create struct {
	currencyRepository *repository.CurrencyRepository
	currencyApi        *service.CurrencyApi
}

func NewCreateHandler(currencyRepository *repository.CurrencyRepository, currencyApi *service.CurrencyApi) *Create {
	return &Create{
		currencyRepository: currencyRepository,
		currencyApi:        currencyApi,
	}
}

func (c *Create) validateCreateRequest(jsonBody map[string]interface{}) *handler.UnprocessableEntity {
	var unprocessableEntity handler.UnprocessableEntity

	symbol, symbolExists := jsonBody["symbol"].(string)
	symbol = strings.Trim(symbol, " ")

	if !symbolExists {
		unprocessableEntity.Add("Symbol is required")
	} else if symbol == "" {
		unprocessableEntity.Add("Symbol can not be empty")
	} else if len([]rune(symbol)) < 3 {
		unprocessableEntity.Add("Symbol should be at least 3 characters long")
	}

	rate, rateExists := jsonBody["rate"].(float64)

	if rateExists && rate < 0 {
		unprocessableEntity.Add("Rate should be equal or greater than zero")
	} else if rate == 0 && !unprocessableEntity.HasErrors() {
		_, currency, err := c.currencyApi.Rate(service.DefaultCurrency, symbol)

		if err != nil {
			unprocessableEntity.Add("Error on validating symbol")
		} else if currency == nil {
			unprocessableEntity.Add("Currency not exists. Rate should be greater than zero")
		}
	}

	return &unprocessableEntity
}

func (c *Create) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var jsonBody map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&jsonBody); err != nil || jsonBody == nil {
		handler.Response(w, http.StatusBadRequest, &handler.ErrorResponse{Message: "Invalid request"})
		return
	}

	errors := c.validateCreateRequest(jsonBody)

	if errors.HasErrors() {
		handler.Response(w, http.StatusUnprocessableEntity, &errors)
		return
	}

	rate, existsRate := jsonBody["rate"].(float64)

	if !existsRate {
		rate = 0
	}

	currency := &entity.Currency{
		Symbol: strings.ToUpper(strings.Trim(jsonBody["symbol"].(string), " ")),
		Rate:   rate,
	}

	err := c.currencyRepository.CreateCurrency(currency)

	if err != nil {
		handler.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	handler.Response(w, http.StatusCreated, &currency)
}
