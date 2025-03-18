package repository

import (
	"database/sql"
	"github.com/Tagliatti/challenge-bravo/entity"
	"strings"
)

type CurrencyRepository struct {
	db *sql.DB
}

func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) FindBySymbol(symbol string) (*entity.Currency, error) {
	rows, err := r.db.Query("SELECT symbol, rate FROM currencies WHERE symbol = UPPER(?)", symbol)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var currency entity.Currency
		if err := rows.Scan(&currency.Symbol, &currency.Rate); err != nil {
			return nil, err
		}
		return &currency, nil
	}

	return nil, nil
}

func (r *CurrencyRepository) FindCurrencySymbols(fromSymbol string, toSymbol string) (from *entity.Currency, to *entity.Currency, err error) {
	rows, err := r.db.Query(`
        SELECT
            symbol,
            rate
        FROM
            currencies
        WHERE
            symbol IN (UPPER(?), UPPER(?))
    `, fromSymbol, toSymbol)

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var symbol string
		var rate float64

		if err = rows.Scan(&symbol, &rate); err != nil {
			return nil, nil, err
		}

		if symbol == strings.ToUpper(fromSymbol) {
			from = &entity.Currency{Symbol: symbol, Rate: rate}
		} else {
			to = &entity.Currency{Symbol: symbol, Rate: rate}
		}
	}

	return from, to, nil
}

func (r *CurrencyRepository) CreateCurrency(currency *entity.Currency) error {
	_, err := r.db.Exec("INSERT INTO currencies (symbol, rate) VALUES (TRIM(UPPER(?)), ?)", currency.Symbol, currency.Rate)

	return err
}

func (r *CurrencyRepository) UpdateCurrency(currency *entity.Currency) error {
	_, err := r.db.Exec("UPDATE currencies SET rate = ? WHERE symbol = UPPER(?)", currency.Rate, currency.Symbol)

	return err
}

func (r *CurrencyRepository) DeleteCurrency(symbol string) (bool, error) {
	result, err := r.db.Exec("DELETE FROM currencies WHERE symbol = UPPER(?)", symbol)

	if err != nil {
		return false, err
	}

	var rows int64 = 0

	if result != nil {
		rows, _ = result.RowsAffected()
	}

	return rows > 0, err
}
