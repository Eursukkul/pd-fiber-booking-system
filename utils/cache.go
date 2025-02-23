package utils

import (
	"errors"
	"sync"

	"github.com/Eursukkul/fiber-booking-system/dto"
)

type (
	// Cache interface
	Cache interface {
		Set(id int, booking *dto.BookingResponse)
		Get(id int) (*dto.BookingResponse, error)
		Delete(id int)
	}

	// InMemoryCache
	InMemoryCache struct {
		store map[int]*dto.BookingResponse
		mu    sync.RWMutex
	}
)

// NewInMemoryCache
func NewInMemoryCache() Cache {
	return &InMemoryCache{
		store: make(map[int]*dto.BookingResponse),
	}
}

// keep data in cache
func (c *InMemoryCache) Set(id int, booking *dto.BookingResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[id] = booking
}

// get data from cache
func (c *InMemoryCache) Get(id int) (*dto.BookingResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	booking, exists := c.store[id]
	if !exists {
		return nil, errors.New("booking not found")
	}
	return booking, nil
}

// delete data from cache	
func (c *InMemoryCache) Delete(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, id)
}
