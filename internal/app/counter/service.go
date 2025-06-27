package counter

import (
	"errors"
	"sync"
)

type PageViewCounter interface {
	Increment(pageID string) error
	GetCount(pageID string) (int, error)
}

type Counter struct {
	mu     sync.RWMutex
	counts map[string]int
	// Поля для хранилища (например, Redis, in-memory map, DB)
}

func New() *Counter {
	mu := sync.RWMutex{}
	mapa := make(map[string]int)

	return &Counter{
		mu:     mu,
		counts: mapa,
	}
}

func (c *Counter) Increment(pageID string) error {
	if pageID == "" {
		return errors.New("pageID cannot be empty")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[pageID]++
	return nil
}

func (c *Counter) GetCount(pageID string) (int, error) {
	if pageID == "" {
		return 0, errors.New("pageID cannot be empty")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	countOfPages, exists := c.counts[pageID]
	if !exists {
		return 0, errors.New("pageID not found")
	}
	return countOfPages, nil
}
