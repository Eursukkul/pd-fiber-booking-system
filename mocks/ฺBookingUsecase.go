// mocks/BookingUsecase.go
package mocks

import (
	"sync"

	"github.com/Eursukkul/fiber-booking-system/dto"
	"github.com/stretchr/testify/mock"
)

// MockBookingUsecase เป็น mock สำหรับ interface BookingUsecase
type MockBookingUsecase struct {
	mock.Mock
}

func (m *MockBookingUsecase) CreateBooking(req dto.BookingRequest) (*dto.BookingResponse, error) {
	args := m.Called(req)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.BookingResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBookingUsecase) GetBookingByID(id int) (*dto.BookingResponse, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.BookingResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBookingUsecase) CancelBooking(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookingUsecase) CheckExpiredBookings() {
	m.Called()
}

func (m *MockBookingUsecase) BackgroundTaskBooking(wg *sync.WaitGroup) {
	m.Called(wg)
}

func (m *MockBookingUsecase) GetAllBookings(sort string, highValue string) ([]*dto.BookingResponse, error) {
	args := m.Called(sort, highValue)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*dto.BookingResponse), args.Error(1)
}

func (m *MockBookingUsecase) UpdateBooking(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockBookingUsecase) UpdateBookingStatus(id int, status string) error {
    args := m.Called(id, status)
    return args.Error(0)
}
