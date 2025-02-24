// MockCache เป็น mock สำหรับ interface Cache

package mocks

import (
	"github.com/Eursukkul/fiber-booking-system/dto"
	"github.com/stretchr/testify/mock"
)
// Cache is a mock for the Cache interface
type Cache struct {
	mock.Mock
}

// Set mocks the Set method of the Cache interface
func (m *Cache) Set(id int, booking *dto.BookingResponse) error {
	args := m.Called(id, booking)
	return args.Error(0)
}

// Get เป็น mock สำหรับเมธอด Get ใน Cache
func (m *Cache) Get(id int) (*dto.BookingResponse, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.BookingResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// Delete เป็น mock สำหรับเมธอด Delete ใน Cache
func (m *Cache) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
