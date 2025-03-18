package service

import (
	"encoding/json"
	"fmt"
	"github.com/Tagliatti/challenge-bravo/entity"
	"io"
	"net/http"
	"strings"
)

const DefaultCurrency = "USD"

const URL = "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/usd.json"

type CurrencyApi struct {
	httpClient *http.Client
}

func NewCurrencyApi(httpClient *http.Client) *CurrencyApi {
	return &CurrencyApi{
		httpClient: httpClient,
	}
}

func (c *CurrencyApi) Rate(fromSymbol string, toSymbol string) (from *entity.Currency, to *entity.Currency, err error) {
	fromSymbol = strings.ToLower(fromSymbol)
	toSymbol = strings.ToLower(toSymbol)

	response, err := c.httpClient.Get(URL)

	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	var jsonResponse map[string]interface{}

	err = json.Unmarshal(body, &jsonResponse)

	if err != nil {
		return nil, nil, err
	}

	if rates, ok := jsonResponse[strings.ToLower(DefaultCurrency)].(map[string]interface{}); ok {
		var rate float64
		if rate, ok = rates[fromSymbol].(float64); ok {
			from = &entity.Currency{
				Symbol: strings.ToUpper(fromSymbol),
				Rate:   rate,
			}
		}
		if rate, ok = rates[toSymbol].(float64); ok {
			to = &entity.Currency{
				Symbol: strings.ToUpper(toSymbol),
				Rate:   rate,
			}
		}

		return from, to, nil
	}

	return nil, nil, fmt.Errorf("rate not found")
}
