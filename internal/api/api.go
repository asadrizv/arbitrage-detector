package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/arbitrage-detector/internal/db"
)

// Server exposes HTTP endpoints.
type Server struct {
	DB *db.DB
}

// Start runs the HTTP server on the given addr.
func (s *Server) Start(addr string) error {
	http.HandleFunc("/opportunities/top", s.handleTop)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleTop(w http.ResponseWriter, r *http.Request) {
	ops := s.DB.TopOpportunities(10, time.Now())
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ops)
}
