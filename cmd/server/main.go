package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/example/arbitrage-detector/internal/api"
	"github.com/example/arbitrage-detector/internal/db"
	"github.com/example/arbitrage-detector/internal/service"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	database, err := db.New("data.json")
	if err != nil {
		log.Fatal(err)
	}

	svc := &service.Service{DB: database}
	svc.StartScheduler()

	server := &api.Server{DB: database}
	log.Println("Starting server on :8080")
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
