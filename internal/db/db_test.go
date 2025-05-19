package db

import (
	"os"
	"testing"
	"time"

	"github.com/example/arbitrage-detector/internal/model"
)

func TestInsertAndTopOpportunities(t *testing.T) {
	tmp, err := os.CreateTemp("", "dbtest-*.json")
	if err != nil {
		t.Fatal(err)
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	defer os.Remove(path)

	d, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	day := time.Now()
	for i := 0; i < 5; i++ {
		o := model.Opportunity{
			CurrencyPair: "EUR/USD",
			Profit:       float64(i),
			CreatedAt:    day,
		}
		if err := d.Insert(o); err != nil {
			t.Fatalf("Insert: %v", err)
		}
	}

	ops := d.TopOpportunities(3, day)
	if len(ops) != 3 {
		t.Fatalf("expected 3 ops, got %d", len(ops))
	}
	if ops[0].Profit < ops[1].Profit {
		t.Errorf("results not sorted desc")
	}
}
