package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/example/arbitrage-detector/internal/db"
	"github.com/example/arbitrage-detector/internal/model"
)

// Service coordinates fetching opportunities and storing them.
type RateProvider interface {
	GetRate(base, quote string) (float64, error)
}

type Service struct {
	DB        *db.DB
	ProviderA RateProvider
	ProviderB RateProvider
	MinSpread float64
}

// defaultProviderA fetches rates from exchangerate.host
type defaultProviderA struct{ client *http.Client }

func (p defaultProviderA) GetRate(base, quote string) (float64, error) {
	url := fmt.Sprintf("https://api.exchangerate.host/latest?base=%s&symbols=%s", base, quote)
	resp, err := p.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var m struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return 0, err
	}
	return m.Rates[quote], nil
}

// defaultProviderB fetches rates from open.er-api.com
type defaultProviderB struct{ client *http.Client }

func (p defaultProviderB) GetRate(base, quote string) (float64, error) {
	url := fmt.Sprintf("https://open.er-api.com/v6/latest/%s", base)
	resp, err := p.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var m struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return 0, err
	}
	return m.Rates[quote], nil
}

// StartScheduler starts a goroutine that fetches data daily.
func (s *Service) StartScheduler() {
	if s.ProviderA == nil {
		s.ProviderA = defaultProviderA{client: http.DefaultClient}
	}
	if s.ProviderB == nil {
		s.ProviderB = defaultProviderB{client: http.DefaultClient}
	}
	if s.MinSpread == 0 {
		s.MinSpread = 0.1
	}
	go func() {
		s.fetchAndStore()
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			s.fetchAndStore()
		}
	}()
}

// fetchAndStore generates fake opportunities and stores them.
func (s *Service) fetchAndStore() {
	pairs := []string{"EUR/USD", "USD/JPY", "GBP/USD", "USD/CHF", "AUD/USD"}
	for _, pair := range pairs {
		parts := strings.Split(pair, "/")
		if len(parts) != 2 {
			continue
		}
		base, quote := parts[0], parts[1]
		r1, err1 := s.ProviderA.GetRate(base, quote)
		r2, err2 := s.ProviderB.GetRate(base, quote)
		if err1 != nil || err2 != nil {
			continue
		}

		if r1 < r2 {
			spread := (r2 - r1) / r1 * 100
			if spread >= s.MinSpread {
				o := model.Opportunity{
					CurrencyPair: pair,
					Profit:       spread,
					CreatedAt:    time.Now(),
				}
				_ = s.DB.Insert(o)
			}
		} else if r2 < r1 {
			spread := (r1 - r2) / r2 * 100
			if spread >= s.MinSpread {
				o := model.Opportunity{
					CurrencyPair: pair,
					Profit:       spread,
					CreatedAt:    time.Now(),
				}
				_ = s.DB.Insert(o)
			}
		}
	}
}
