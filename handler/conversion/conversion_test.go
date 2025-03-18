package conversion

import (
	"bytes"
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

func TestConversionValidations(t *testing.T) {
	testCases := []test.ApiTestCase{
		{"Should return validation error (Missing 'from' parameter)", `/?to=USD&amount=1`, http.StatusUnprocessableEntity, "Missing 'from' parameter"},
		{"Should return validation error (Missing 'to' parameter)", `/?from=USD&amount=1`, http.StatusUnprocessableEntity, "Missing 'to' parameter"},
		{"Should return validation error (Missing 'amount' parameter)", `/?from=USD&to=EUR`, http.StatusUnprocessableEntity, "Missing 'amount' parameter"},
		{"Should return validation error (Invalid 'amount' parameter)", `/?from=USD&to=EUR&amount=ABC`, http.StatusUnprocessableEntity, "Invalid 'amount' parameter"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db := test.SetUp(t)
			currencyRepository := repository.NewCurrencyRepository(db)
			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tc.Content, nil)
			httpClient := test.NewHttpClientWithError(t)
			currencyApi := service.NewCurrencyApi(httpClient)
			handler := NewConversionHandler(currencyRepository, currencyApi)
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

func TestConversionValidationCurrencyNotFound(t *testing.T) {
	testCases := []test.ApiTestCase{
		{"Should return currency not found (from)", `/?from=ABC&to=EUR&amount=1`, http.StatusBadRequest, `{"message":"currency 'ABC' or 'EUR' not found"}`},
		{"Should return currency not found (to)", `/?from=EUR&to=ABC&amount=1`, http.StatusBadRequest, `{"message":"currency 'EUR' or 'ABC' not found"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db := test.SetUp(t)
			currencyRepository := repository.NewCurrencyRepository(db)
			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tc.Content, nil)
			httpClient := test.NewHttpClientWithError(t)
			currencyApi := service.NewCurrencyApi(httpClient)
			handler := NewConversionHandler(currencyRepository, currencyApi)
			handler.Handler(responseRecorder, request)

			if responseRecorder.Code != tc.ExpectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, tc.ExpectedCode)
			}
			if strings.Trim(responseRecorder.Body.String(), "\n") != tc.ExpectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), tc.ExpectedBody)
			}
		})
	}
}

func TestConversionLocal(t *testing.T) {
	t.Run("Should return success (from offline, to offline)", func(t *testing.T) {
		db := test.SetUp(t)
		currencyRepository := repository.NewCurrencyRepository(db)

		from := &entity.Currency{Symbol: "CUR_FROM", Rate: 1}
		to := &entity.Currency{Symbol: "CUR_TO", Rate: 2}

		err := currencyRepository.CreateCurrency(from)

		if err != nil {
			t.Errorf("Error creating currency from: %v", err)
		}

		err = currencyRepository.CreateCurrency(to)

		if err != nil {
			t.Errorf("Error creating currency to: %v", err)
		}

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/?from=CUR_FROM&to=CUR_TO&amount=1", nil)
		httpClient := test.NewHttpClientWithError(t)
		currencyApi := service.NewCurrencyApi(httpClient)
		handler := NewConversionHandler(currencyRepository, currencyApi)
		handler.Handler(responseRecorder, request)

		expectedCode := http.StatusOK
		expectedBody := `{"from":"CUR_FROM","to":"CUR_TO","amount":1,"result":2}`

		if responseRecorder.Code != expectedCode {
			t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, expectedCode)
		}
		if strings.Trim(responseRecorder.Body.String(), "\n") != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expectedBody)
		}
	})
}

func TestConversionOnlineSymbolsNotFound(t *testing.T) {
	t.Run("Should return error (rate not found)", func(t *testing.T) {
		db := test.SetUp(t)
		currencyRepository := repository.NewCurrencyRepository(db)
		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/?from=USD&to=EUR&amount=1", nil)
		httpClient := &http.Client{Transport: test.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			}, nil
		})}
		currencyApi := service.NewCurrencyApi(httpClient)
		handler := NewConversionHandler(currencyRepository, currencyApi)
		handler.Handler(responseRecorder, request)

		expectedCode := http.StatusInternalServerError
		expectedBody := `{"message":"rate not found"}`

		if responseRecorder.Code != expectedCode {
			t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, expectedCode)
		}
		if strings.Trim(responseRecorder.Body.String(), "\n") != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), expectedBody)
		}
	})
}

func TestConversion(t *testing.T) {
	testCases := []test.ApiTestCase{
		{"Should return success (from online, to online)", "/?from=USD&to=EUR&amount=1", http.StatusOK, `{"from":"USD","to":"EUR","amount":1,"result":2}`},
		{"Should return success (from online, to offline)", "/?from=USD&to=LOCAL&amount=1", http.StatusOK, `{"from":"USD","to":"LOCAL","amount":1,"result":3}`},
		{"Should return success (from offline, to online)", "/?from=LOCAL&to=USD&amount=3", http.StatusOK, `{"from":"LOCAL","to":"USD","amount":3,"result":1}`},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db := test.SetUp(t)
			currencyRepository := repository.NewCurrencyRepository(db)

			currency := &entity.Currency{Symbol: "LOCAL", Rate: 3}
			err := currencyRepository.CreateCurrency(currency)

			if err != nil {
				t.Errorf("Error creating currency: %v", err)
			}

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tc.Content, nil)
			httpClient := &http.Client{Transport: test.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"usd":{"usd":1,"eur":2}}`)),
				}, nil
			})}
			currencyApi := service.NewCurrencyApi(httpClient)
			handler := NewConversionHandler(currencyRepository, currencyApi)
			handler.Handler(responseRecorder, request)

			if responseRecorder.Code != tc.ExpectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", responseRecorder.Code, tc.ExpectedCode)
			}
			if strings.Trim(responseRecorder.Body.String(), "\n") != tc.ExpectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", responseRecorder.Body.String(), tc.ExpectedBody)
			}
		})
	}
}
