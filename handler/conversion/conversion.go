package conversion

import (
	"fmt"
	"github.com/Tagliatti/challenge-bravo/entity"
	"github.com/Tagliatti/challenge-bravo/handler"
	"github.com/Tagliatti/challenge-bravo/repository"
	"github.com/Tagliatti/challenge-bravo/service"
	"net/http"
	"strconv"
)

type Conversion struct {
	currencyRepository *repository.CurrencyRepository
	currencyApi        *service.CurrencyApi
}

func NewConversionHandler(currencyRepository *repository.CurrencyRepository, currencyApi *service.CurrencyApi) *Conversion {
	return &Conversion{
		currencyRepository: currencyRepository,
		currencyApi:        currencyApi,
	}
}

type convertResponse struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Result float64 `json:"result"`
}

func validateConversionRequest(r *http.Request) *handler.UnprocessableEntity {
	var unprocessableEntity handler.UnprocessableEntity

	if r.URL.Query().Get("from") == "" {
		unprocessableEntity.Add("Missing 'from' parameter")
	}
	if r.URL.Query().Get("to") == "" {
		unprocessableEntity.Add("Missing 'to' parameter")
	}
	if r.URL.Query().Get("amount") == "" {
		unprocessableEntity.Add("Missing 'amount' parameter")
	} else if _, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64); err != nil {
		unprocessableEntity.Add("Invalid 'amount' parameter")
	}

	return &unprocessableEntity
}

func convert(from *entity.Currency, to *entity.Currency, amount float64) float64 {
	return (amount / from.Rate) * to.Rate
}

func (c *Conversion) Handler(w http.ResponseWriter, r *http.Request) {
	errors := validateConversionRequest(r)
	if errors.HasErrors() {
		handler.Response(w, http.StatusUnprocessableEntity, &errors)
		return
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amount, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)

	fromLocalCurrency, toLocalCurrency, err := c.currencyRepository.FindCurrencySymbols(from, to)

	if err != nil {
		handler.ResponseError(w, http.StatusInternalServerError, err)
		return
	}
	if fromLocalCurrency == nil || toLocalCurrency == nil {
		handler.Response(w, http.StatusBadRequest, &handler.ErrorResponse{Message: fmt.Sprintf("currency '%s' or '%s' not found", from, to)})
		return
	}

	var result float64

	if fromLocalCurrency.Rate > 0 && toLocalCurrency.Rate > 0 {
		result = convert(fromLocalCurrency, toLocalCurrency, amount)
	} else {
		fromOnlineCurrency, toOnlineCurrency, err := c.currencyApi.Rate(from, to)

		if err != nil {
			handler.ResponseError(w, http.StatusInternalServerError, err)
			return
		}

		if fromOnlineCurrency != nil && toOnlineCurrency != nil {
			result = convert(fromOnlineCurrency, toOnlineCurrency, amount)
		} else if fromLocalCurrency.Rate > 0 && toOnlineCurrency != nil {
			result = convert(fromLocalCurrency, toOnlineCurrency, amount)
		} else if fromOnlineCurrency != nil && toLocalCurrency.Rate > 0 {
			result = convert(fromOnlineCurrency, toLocalCurrency, amount)
		} else {
			handler.Response(w, http.StatusBadRequest, &handler.ErrorResponse{Message: "conversion not available"})
			return
		}
	}

	convertResponse := convertResponse{
		From:   from,
		To:     to,
		Amount: amount,
		Result: result,
	}

	handler.Response(w, http.StatusOK, &convertResponse)
}
