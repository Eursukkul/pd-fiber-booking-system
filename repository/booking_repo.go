package repository

import (
	"sync"
	"time"

	"github.com/Eursukkul/fiber-booking-system/dto"

)

type (
	BookingRepository interface{
		Create(req dto.BookingRequest) *dto.BookingResponse
		GetByID(id int) (*dto.BookingResponse, bool)
		GetAll() []*dto.BookingResponse
		Update(id int, status string) bool
	}

	MockBookingRepository struct {
		bookings map[int]dto.BookingResponse
		mu       sync.RWMutex
	}
)

func NewMockBookingRepository() BookingRepository {
	booking := make(map[int]dto.BookingResponse)
	loc, _ := time.LoadLocation("Asia/Bangkok")
	for i := 1; i <= 10; i++ {
		booking[i] = dto.BookingResponse{
			ID: i,
			UserID: i,
			ServiceID: i,
			Price: float64(i * 1000),
			Status: "pending",
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Minute).In(loc).Format(time.RFC3339),
			UpdatedAt: time.Now().In(loc).Format(time.RFC3339),
		}
	}
	return &MockBookingRepository{
		bookings: booking,
		mu:       sync.RWMutex{},
	}
}

// Create
func (m *MockBookingRepository) Create(req dto.BookingRequest) *dto.BookingResponse {
    m.mu.Lock()
    defer m.mu.Unlock()
    id := len(m.bookings) + 1
    booking := &dto.BookingResponse{
        ID:        id,
        UserID:    req.UserID,
        ServiceID: req.ServiceID,
        Price:     req.Price,
        Status:    "pending",
        CreatedAt: time.Now().Format(time.RFC3339),
        UpdatedAt: time.Now().Format(time.RFC3339),
    }
    m.bookings[id] = *booking
    return booking
}

// GetByID
func (m *MockBookingRepository) GetByID(id int) (*dto.BookingResponse, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    booking, exists := m.bookings[id]
    	return &booking, exists
}

// GetAll
func (m *MockBookingRepository) GetAll() []*dto.BookingResponse {
	m.mu.RLock()
	defer m.mu.RUnlock()
	bookings := []*dto.BookingResponse{}
	for _, b := range m.bookings {
		bookingCopy := b // create copy to avoid pointer issue
		bookings = append(bookings, &bookingCopy)
	}
	return bookings
}

// Update status booking
func (m *MockBookingRepository) Update(id int, status string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	booking, exists := m.bookings[id]
	if !exists {
		return false
	}
	booking.Status = status
	booking.UpdatedAt = time.Now().Format(time.RFC3339)
	m.bookings[id] = booking
	return true
}

