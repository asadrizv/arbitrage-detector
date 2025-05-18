package db

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/example/arbitrage-detector/internal/model"
)

// DB is a simple file-based storage.
type DB struct {
	filePath string
	mu       sync.Mutex
	data     []model.Opportunity
	nextID   int
}

// New creates a new database backed by the given file.
func New(file string) (*DB, error) {
	d := &DB{filePath: file}
	if err := d.load(); err != nil {
		return nil, err
	}
	return d, nil
}

// load reads opportunities from file if it exists.
func (d *DB) load() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	f, err := os.Open(d.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			d.data = []model.Opportunity{}
			d.nextID = 1
			return nil
		}
		return err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&d.data); err != nil {
		return err
	}
	// set nextID
	maxID := 0
	for _, o := range d.data {
		if o.ID > maxID {
			maxID = o.ID
		}
	}
	d.nextID = maxID + 1
	return nil
}

// save writes opportunities to file.
func (d *DB) save() error {
	f, err := os.Create(d.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(d.data)
}

// Insert adds a new opportunity to the DB.
func (d *DB) Insert(o model.Opportunity) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	o.ID = d.nextID
	d.nextID++
	d.data = append(d.data, o)
	return d.save()
}

// TopOpportunities returns top n opportunities for the given date (local time).
func (d *DB) TopOpportunities(n int, day time.Time) []model.Opportunity {
	d.mu.Lock()
	defer d.mu.Unlock()

	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)

	var todays []model.Opportunity
	for _, o := range d.data {
		if !o.CreatedAt.Before(start) && o.CreatedAt.Before(end) {
			todays = append(todays, o)
		}
	}
	sort.Slice(todays, func(i, j int) bool { return todays[i].Profit > todays[j].Profit })
	if len(todays) > n {
		todays = todays[:n]
	}
	result := make([]model.Opportunity, len(todays))
	copy(result, todays)
	return result
}
