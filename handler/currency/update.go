package currency

import (
	"encoding/json"
	"github.com/Tagliatti/challenge-bravo/handler"
	"github.com/Tagliatti/challenge-bravo/repository"
	"github.com/Tagliatti/challenge-bravo/service"
	"net/http"
	"strings"
)

type Update struct {
	currencyRepository *repository.CurrencyRepository
	currencyApi        *service.CurrencyApi
}

func NewUpdateHandler(currencyRepository *repository.CurrencyRepository, currencyApi *service.CurrencyApi) *Update {
	return &Update{
		currencyRepository: currencyRepository,
		currencyApi:        currencyApi,
	}
}

func (u *Update) validateUpdateRequest(symbol string, jsonBody map[string]interface{}) *handler.UnprocessableEntity {
	var unprocessableEntity handler.UnprocessableEntity

	rate, rateExists := jsonBody["rate"].(float64)

	if !rateExists {
		unprocessableEntity.Add("Rate is required")
	} else if rate < 0 {
		unprocessableEntity.Add("Rate should be equal or greater than zero")
	} else if rate == 0 && !unprocessableEntity.HasErrors() {
		_, currency, err := u.currencyApi.Rate(service.DefaultCurrency, symbol)

		if err != nil {
			unprocessableEntity.Add("Error on validating symbol")
		} else if currency == nil {
			unprocessableEntity.Add("Currency not exists. Rate should be greater than zero")
		}
	}

	return &unprocessableEntity
}

func (u *Update) Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	symbol := strings.Trim(r.PathValue("symbol"), " ")
	currency, err := u.currencyRepository.FindBySymbol(symbol)

	if err != nil {
		handler.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	if currency == nil {
		handler.Response(w, http.StatusNotFound, &handler.ErrorResponse{Message: "currency not found"})
		return
	}

	var jsonBody map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&jsonBody); err != nil || jsonBody == nil {
		handler.Response(w, http.StatusBadRequest, &handler.ErrorResponse{Message: "invalid request"})
		return
	}

	errors := u.validateUpdateRequest(symbol, jsonBody)

	if errors.HasErrors() {
		handler.Response(w, http.StatusUnprocessableEntity, &errors)
		return
	}

	rate, _ := jsonBody["rate"].(float64)

	currency.Rate = rate

	err = u.currencyRepository.UpdateCurrency(currency)

	if err != nil {
		handler.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	handler.Response(w, http.StatusOK, &currency)
}
