package service

import (
	"math/rand"
	"time"

	"github.com/example/arbitrage-detector/internal/db"
	"github.com/example/arbitrage-detector/internal/model"
)

// Service coordinates fetching opportunities and storing them.
type Service struct {
	DB *db.DB
}

// StartScheduler starts a goroutine that fetches data daily.
func (s *Service) StartScheduler() {
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
	for i := 0; i < 20; i++ {
		o := model.Opportunity{
			CurrencyPair: pairs[rand.Intn(len(pairs))],
			Profit:       rand.Float64() * 5,
			CreatedAt:    time.Now(),
		}
		_ = s.DB.Insert(o)
	}
}
