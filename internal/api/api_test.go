package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/example/arbitrage-detector/internal/db"
	"github.com/example/arbitrage-detector/internal/model"
)

func TestHandleTop(t *testing.T) {
	tmp, err := os.CreateTemp("", "apidb-*.json")
	if err != nil {
		t.Fatal(err)
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	defer os.Remove(path)

	d, err := db.New(path)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	for i := 0; i < 3; i++ {
		_ = d.Insert(model.Opportunity{CurrencyPair: "EUR/USD", Profit: float64(10 - i), CreatedAt: now})
	}

	s := &Server{DB: d}
	req := httptest.NewRequest(http.MethodGet, "/opportunities/top", nil)
	w := httptest.NewRecorder()
	s.handleTop(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}

	var ops []model.Opportunity
	if err := json.Unmarshal(w.Body.Bytes(), &ops); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(ops) != 3 {
		t.Fatalf("expected 3 ops, got %d", len(ops))
	}
}
