package service

import (
	"os"
	"testing"
	"time"

	"github.com/example/arbitrage-detector/internal/db"
)

func TestFetchAndStore(t *testing.T) {
	tmp, err := os.CreateTemp("", "svcdb-*.json")
	if err != nil {
		t.Fatal(err)
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	defer os.Remove(path)

	d, err := db.New(path)
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}

	pA := stubProvider{rate: 1.0}
	pB := stubProvider{rate: 1.1}
	svc := &Service{DB: d, ProviderA: pA, ProviderB: pB}
	svc.fetchAndStore()

	if len(d.TopOpportunities(100, time.Now())) == 0 {
		t.Error("expected opportunities inserted")
	}
}

type stubProvider struct {
	rate float64
	err  error
}

func (s stubProvider) GetRate(base, quote string) (float64, error) {
	return s.rate, s.err
}
