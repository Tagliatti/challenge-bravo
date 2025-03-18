package currency

import (
	"github.com/Tagliatti/challenge-bravo/handler"
	"github.com/Tagliatti/challenge-bravo/repository"
	"net/http"
	"strings"
)

type Delete struct {
	currencyRepository *repository.CurrencyRepository
}

func NewDeleteHandler(currencyRepository *repository.CurrencyRepository) *Delete {
	return &Delete{
		currencyRepository: currencyRepository,
	}
}

func (d *Delete) Handler(w http.ResponseWriter, r *http.Request) {
	found, err := d.currencyRepository.DeleteCurrency(strings.Trim(r.PathValue("symbol"), " "))

	if err != nil {
		handler.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	if found {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
