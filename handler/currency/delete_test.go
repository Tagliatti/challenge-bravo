package currency

import (
	"github.com/Tagliatti/challenge-bravo/entity"
	"github.com/Tagliatti/challenge-bravo/repository"
	"github.com/Tagliatti/challenge-bravo/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNetFoundErrorWithoutPassingSymbol(t *testing.T) {
	t.Run("Should return 404 status code (without passing symbol)", func(t *testing.T) {
		db := test.SetUp(t)
		currencyRepository := repository.NewCurrencyRepository(db)
		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/", nil)
		handler := NewDeleteHandler(currencyRepository)
		handler.Handler(responseRecorder, request)

		if responseRecorder.Code != http.StatusNotFound {
			t.Errorf("Expecting HTTP status code %d, got %d", http.StatusNotFound, responseRecorder.Code)
		}
	})
}

func TestNetFoundErrorWithPassingSymbol(t *testing.T) {
	t.Run("Should return 404 status code (with passing symbol)", func(t *testing.T) {
		db := test.SetUp(t)
		currencyRepository := repository.NewCurrencyRepository(db)
		err := currencyRepository.CreateCurrency(&entity.Currency{Symbol: "ABC", Rate: 1})

		if err != nil {
			t.Errorf("Error on create currency: %s", err.Error())
		}

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/NOT_FOUND", nil)
		handler := NewDeleteHandler(currencyRepository)
		handler.Handler(responseRecorder, request)

		if responseRecorder.Code != http.StatusNotFound {
			t.Errorf("Expecting HTTP status code %d, got %d", http.StatusNotFound, responseRecorder.Code)
		}
	})
}

func TestDeleteSymbol(t *testing.T) {
	testCases := []test.ApiTestCase{
		{"Should delete symbol and return 204 status code (symbol uppercase)", `ABC`, http.StatusNoContent, ``},
		{"Should delete symbol and return 204 status code (symbol lowercase)", `abc`, http.StatusNoContent, ``},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db := test.SetUp(t)
			currencyRepository := repository.NewCurrencyRepository(db)
			err := currencyRepository.CreateCurrency(&entity.Currency{Symbol: tc.Content, Rate: 1})

			if err != nil {
				t.Errorf("Error on create currency: %s", err.Error())
			}

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodDelete, "/"+tc.Content, nil)
			request.SetPathValue("symbol", tc.Content)
			handler := NewDeleteHandler(currencyRepository)
			handler.Handler(responseRecorder, request)

			if responseRecorder.Code != tc.ExpectedCode {
				t.Errorf("Expecting HTTP status code %d, got %d", tc.ExpectedCode, responseRecorder.Code)
			}

			if responseRecorder.Body.String() != tc.ExpectedBody {
				t.Errorf("Expecting body %s, got %s", tc.ExpectedBody, responseRecorder.Body.String())
			}

			currency, err := currencyRepository.FindBySymbol(tc.Content)

			if err != nil {
				t.Errorf("Error on find by symbol: %s", err.Error())
			}

			if currency != nil {
				t.Errorf("Currency shouldn't be found")
			}
		})
	}
}
