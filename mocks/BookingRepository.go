// MockBookingRepository เป็น mock สำหรับ interface BookingRepository

package mocks

import (
	"github.com/Eursukkul/fiber-booking-system/dto"
	"github.com/stretchr/testify/mock"
)

// MockBookingRepository 
type MockBookingRepository struct {
	mock.Mock
}

// Create mock data
func (m *MockBookingRepository) Create(req dto.BookingRequest) *dto.BookingResponse {
	args := m.Called(req)
	return args.Get(0).(*dto.BookingResponse)
}

// GetByID mock data
func (m *MockBookingRepository) GetByID(id int) (*dto.BookingResponse, bool) {
	args := m.Called(id)
	return args.Get(0).(*dto.BookingResponse), args.Bool(1)
}

// GetAll mock data
func (m *MockBookingRepository) GetAll() []*dto.BookingResponse {
	args := m.Called()
	return args.Get(0).([]*dto.BookingResponse)
}

// GetHighValueBookings mock data
func (m *MockBookingRepository) GetHighValueBookings(threshold float64) []*dto.BookingResponse {
	args := m.Called(threshold)
	return args.Get(0).([]*dto.BookingResponse)
}

// UpdateBookingStatus mock data
func (m *MockBookingRepository) UpdateBookingStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}
