package entity

type Currency struct {
	Symbol string  `json:"symbol"`
	Rate   float64 `json:"rate"`
}
