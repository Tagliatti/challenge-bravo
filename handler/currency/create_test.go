package currency

import (
	"bytes"
	"errors"
	"github.com/Tagliatti/challenge-bravo/repository"
	"github.com/Tagliatti/challenge-bravo/service"
	"github.com/Tagliatti/challenge-bravo/test"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateInvalidRequest(t *testing.T) {
	db := test.SetUp(t)
	currencyRepository := repository.NewCurrencyRepository(db)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	httpClient := test.NewHttpClientWithError(t)
	currencyApi := service.NewCurrencyApi(httpClient)
	handler := NewCreateHandler(currencyRepository, currencyApi)
	handler.Handler(responseRecorder, request)

	expectedCode := http.StatusBadRequest
	expectedBody := "Invalid request"

	if responseRecorder.Code != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, expectedCode)
	}
	if !strings.Contains(responseRecorder.Body.String(), expectedBody) {
		t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expectedBody)
	}
}

func TestCreateValidations(t *testing.T) {
	testCases := []test.ApiTestCase{
		{"Should return validation error (Symbol is required)", `{}`, http.StatusUnprocessableEntity, "Symbol is required"},
		{"Should return validation error (Symbol can not be empty)", `{"symbol": ""}`, http.StatusUnprocessableEntity, "Symbol can not be empty"},
		{"Should return validation error (Symbol should be at least 3 characters long)", `{"symbol": "AA"}`, http.StatusUnprocessableEntity, "Symbol should be at least 3 characters long"},
		{"Should return validation error (Rate should be equal or greater than zero)", `{"rate": -1}`, http.StatusUnprocessableEntity, "Rate should be equal or greater than zero"},
	}

	httpClient := test.NewHttpClientWithError(t)
	currencyApi := service.NewCurrencyApi(httpClient)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db := test.SetUp(t)
			currencyRepository := repository.NewCurrencyRepository(db)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tc.Content))

			handler := NewCreateHandler(currencyRepository, currencyApi)
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

func TestCreateValidationRate0(t *testing.T) {
	testCases := []test.ApiTestCase{
		{"Should return validation error (Error on validating symbol)", `{"symbol": "ERROR", "rate": 0}`, http.StatusUnprocessableEntity, "Error on validating symbol"},
		{"Should return validation error (Currency not exists. Rate should be greater than zero)", `{"symbol": "SUCCESS", "rate": 0}`, http.StatusUnprocessableEntity, "Currency not exists. Rate should be greater than zero"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db := test.SetUp(t)
			currencyRepository := repository.NewCurrencyRepository(db)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tc.Content))
			httpClient := &http.Client{Transport: test.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				if strings.Contains(tc.Content, "ERROR") {
					return nil, errors.New("internal server error")
				} else {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{"usd": {}}`)),
					}, nil
				}
			})}
			currencyApi := service.NewCurrencyApi(httpClient)
			handler := NewCreateHandler(currencyRepository, currencyApi)
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

func TestCreateWithRateGreaterThan0(t *testing.T) {
	t.Run("Should create currency with rate greater than zero", func(t *testing.T) {
		db := test.SetUp(t)
		currencyRepository := repository.NewCurrencyRepository(db)

		expectedBody := `{"symbol":"ABC","rate":1.245}`
		expectedCode := http.StatusCreated

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(expectedBody))
		httpClient := test.NewHttpClientWithError(t)
		currencyApi := service.NewCurrencyApi(httpClient)
		handler := NewCreateHandler(currencyRepository, currencyApi)
		handler.Handler(responseRecorder, request)

		if responseRecorder.Code != expectedCode {
			t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, expectedCode)
		}
		if strings.Trim(responseRecorder.Body.String(), "\n") != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expectedBody)
		}

		currency, err := currencyRepository.FindBySymbol("ABC")

		if err != nil {
			t.Errorf("Error on find by symbol: %s", err.Error())
		}

		if currency == nil {
			t.Errorf("Currency should be found")
		}
	})
}

func TestCreateWithRate0(t *testing.T) {
	expectedBody := `{"symbol":"ABC","rate":0}`

	testCases := []test.ApiTestCase{
		{"Should create currency with rate zero (passing rate 0)", `{"symbol":"ABC","rate":0}`, http.StatusCreated, expectedBody},
		{"Should create currency with rate zero (not passing field rate)", `{"symbol":"ABC"}`, http.StatusCreated, expectedBody},
		{"Should create currency without space in symbol", `{"symbol":" ABC "}`, http.StatusCreated, expectedBody},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db := test.SetUp(t)
			currencyRepository := repository.NewCurrencyRepository(db)

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tc.Content))
			httpClient := &http.Client{Transport: test.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"usd": {"abc": 1}}`)),
				}, nil
			})}
			currencyApi := service.NewCurrencyApi(httpClient)
			handler := NewCreateHandler(currencyRepository, currencyApi)
			handler.Handler(responseRecorder, request)

			if responseRecorder.Code != tc.ExpectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, tc.ExpectedCode)
			}
			if strings.Trim(responseRecorder.Body.String(), "\n") != tc.ExpectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), tc.ExpectedBody)
			}

			currency, err := currencyRepository.FindBySymbol("ABC")

			if err != nil {
				t.Errorf("Error on find by symbol: %s", err.Error())
			}

			if currency == nil {
				t.Errorf("Currency should be found")
			}
		})
	}
}
