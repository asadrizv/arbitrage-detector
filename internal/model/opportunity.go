package model

import "time"

// Opportunity represents a forex arbitrage opportunity.
type Opportunity struct {
	ID           int       `json:"id"`
	CurrencyPair string    `json:"currency_pair"`
	Profit       float64   `json:"profit"`
	CreatedAt    time.Time `json:"created_at"`
}
