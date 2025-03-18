package currency

import (
	"bytes"
	"errors"
	"github.com/Tagliatti/challenge-bravo/entity"
	"github.com/Tagliatti/challenge-bravo/repository"
	"github.com/Tagliatti/challenge-bravo/service"
	"github.com/Tagliatti/challenge-bravo/test"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateInvalidRequest(t *testing.T) {
	db := test.SetUp(t)
	symbol := "ABC"
	currencyRepository := repository.NewCurrencyRepository(db)

	currency := &entity.Currency{Symbol: symbol, Rate: 0}
	currencyRepository.CreateCurrency(currency)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/"+symbol, nil)
	request.SetPathValue("symbol", symbol)
	httpClient := test.NewHttpClientWithError(t)
	currencyApi := service.NewCurrencyApi(httpClient)
	handler := NewUpdateHandler(currencyRepository, currencyApi)
	handler.Handler(responseRecorder, request)

	expectedCode := http.StatusBadRequest
	expectedBody := "invalid request"

	if responseRecorder.Code != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, expectedCode)
	}
	if !strings.Contains(responseRecorder.Body.String(), expectedBody) {
		t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expectedBody)
	}
}

func TestUpdateValidations(t *testing.T) {
	testCases := []test.ApiTestCase{
		{"Should return validation error (Rate is required)", `{}`, http.StatusUnprocessableEntity, "Rate is required"},
		{"Should return validation error (Rate should be equal or greater than zero)", `{"rate": -1}`, http.StatusUnprocessableEntity, "Rate should be equal or greater than zero"},
	}

	httpClient := test.NewHttpClientWithError(t)
	currencyApi := service.NewCurrencyApi(httpClient)
	symbol := "ABC"

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db := test.SetUp(t)
			currencyRepository := repository.NewCurrencyRepository(db)

			currency := &entity.Currency{Symbol: symbol, Rate: 0}
			currencyRepository.CreateCurrency(currency)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString(tc.Content))
			request.SetPathValue("symbol", symbol)

			handler := NewUpdateHandler(currencyRepository, currencyApi)
			handler.Handler(responseRecorder, request)

			if responseRecorder.Code != tc.ExpectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, tc.ExpectedCode)
			}
			if !strings.Contains(responseRecorder.Body.String(), tc.ExpectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), tc.ExpectedBody)
			}
		})
	}
}

func TestUpdateValidationRate0(t *testing.T) {
	testCases := []test.ApiTestCase{
		{"Should return validation error (Error on validating symbol)", `{"rate": 0}`, http.StatusUnprocessableEntity, "Error on validating symbol"},
		{"Should return validation error (Currency not exists. Rate should be greater than zero)", `{"rate": 0}`, http.StatusUnprocessableEntity, "Currency not exists. Rate should be greater than zero"},
	}

	symbol := "ABC"

	for index, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db := test.SetUp(t)
			currencyRepository := repository.NewCurrencyRepository(db)

			currency := &entity.Currency{Symbol: symbol, Rate: 0}
			currencyRepository.CreateCurrency(currency)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPut, "/"+symbol, bytes.NewBufferString(tc.Content))
			request.SetPathValue("symbol", symbol)
			httpClient := &http.Client{Transport: test.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				if index == 0 {
					return nil, errors.New("internal server error")
				} else {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{"usd": {}}`)),
					}, nil
				}
			})}
			currencyApi := service.NewCurrencyApi(httpClient)
			handler := NewUpdateHandler(currencyRepository, currencyApi)
			handler.Handler(responseRecorder, request)

			if responseRecorder.Code != tc.ExpectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, tc.ExpectedCode)
			}
			if !strings.Contains(responseRecorder.Body.String(), tc.ExpectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), tc.ExpectedBody)
			}
		})
	}
}

func TestUpdateNotFound(t *testing.T) {
	t.Run("Should return not found", func(t *testing.T) {
		db := test.SetUp(t)
		symbol := "ABC"
		currencyRepository := repository.NewCurrencyRepository(db)
		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/"+symbol, nil)
		request.SetPathValue("symbol", symbol)
		httpClient := test.NewHttpClientWithError(t)
		currencyApi := service.NewCurrencyApi(httpClient)
		handler := NewUpdateHandler(currencyRepository, currencyApi)
		handler.Handler(responseRecorder, request)

		expectedBody := `{"message":"currency not found"}`

		if responseRecorder.Code != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, http.StatusNotFound)
		}

		if strings.Trim(responseRecorder.Body.String(), "\n") != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expectedBody)
		}
	})
}

func TestUpdateWithRateGreaterThan0(t *testing.T) {
	t.Run("Should create currency with rate greater than zero", func(t *testing.T) {
		db := test.SetUp(t)
		currencyRepository := repository.NewCurrencyRepository(db)

		symbol := "ABC"
		expectedBody := `{"symbol":"ABC","rate":1.245}`
		requestBody := `{"rate":1.245}`
		expectedCode := http.StatusOK
		expectedRate := 1.245

		currency := &entity.Currency{Symbol: symbol, Rate: 0}
		currencyRepository.CreateCurrency(currency)

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/"+symbol, bytes.NewBufferString(requestBody))
		request.SetPathValue("symbol", symbol)
		httpClient := test.NewHttpClientWithError(t)
		currencyApi := service.NewCurrencyApi(httpClient)
		handler := NewUpdateHandler(currencyRepository, currencyApi)
		handler.Handler(responseRecorder, request)

		if responseRecorder.Code != expectedCode {
			t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, expectedCode)
		}
		if strings.Trim(responseRecorder.Body.String(), "\n") != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expectedBody)
		}

		currency, err := currencyRepository.FindBySymbol(symbol)

		if err != nil {
			t.Errorf("error on find by symbol: %s", err.Error())
		}

		if currency == nil {
			t.Errorf("currency should be found")
		}
		if currency.Rate != expectedRate {
			t.Errorf("currency rate: got %f want %f", currency.Rate, expectedRate)
		}
	})
}

func TestUpdateWithRate0(t *testing.T) {
	t.Run("Should create currency with rate zero", func(t *testing.T) {
		db := test.SetUp(t)
		currencyRepository := repository.NewCurrencyRepository(db)

		symbol := "ABC"
		expectedBody := `{"symbol":"ABC","rate":0}`
		requestBody := `{"rate":0}`
		expectedCode := http.StatusOK
		expectedRate := 0.0

		currency := &entity.Currency{Symbol: symbol, Rate: 1.234}
		currencyRepository.CreateCurrency(currency)

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/"+symbol, bytes.NewBufferString(requestBody))
		request.SetPathValue("symbol", symbol)
		httpClient := &http.Client{Transport: test.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"usd": {"abc": 1}}`)),
			}, nil
		})}
		currencyApi := service.NewCurrencyApi(httpClient)
		handler := NewUpdateHandler(currencyRepository, currencyApi)
		handler.Handler(responseRecorder, request)

		if responseRecorder.Code != expectedCode {
			t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, expectedCode)
		}
		if strings.Trim(responseRecorder.Body.String(), "\n") != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expectedBody)
		}

		currency, err := currencyRepository.FindBySymbol(symbol)

		if err != nil {
			t.Errorf("error on find by symbol: %s", err.Error())
		}
		if currency == nil {
			t.Errorf("currency should be found")
		}
		if currency.Rate != expectedRate {
			t.Errorf("currency rate: got %f want %f", currency.Rate, expectedRate)
		}
	})
}
