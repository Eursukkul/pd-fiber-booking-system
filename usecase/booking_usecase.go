package usecase

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Eursukkul/fiber-booking-system/dto"
	"github.com/Eursukkul/fiber-booking-system/repository"
	"github.com/Eursukkul/fiber-booking-system/utils"
)

type (
	BookingUsecase interface {
		CreateBooking(req dto.BookingRequest) (*dto.BookingResponse, error)
		GetBookingByID(id int) (*dto.BookingResponse, error)
		GetAllBookings(sortBy string, highValue string) ([]*dto.BookingResponse, error)
		UpdateBooking(id int, status string) error
		CancelBooking(id int) error
		BackgroundTaskBooking(wg *sync.WaitGroup)
		UpdateBookingStatus(id int, status string) error
	}

	bookingUsecase struct {
		repo  repository.BookingRepository
		cache utils.Cache
		mu    sync.RWMutex
	}
)

func NewBookingUsecase(repo repository.BookingRepository, cache utils.Cache) BookingUsecase {
	return &bookingUsecase{
		repo:  repo,
		cache: cache,
	}
}

// Create
func (u *bookingUsecase) CreateBooking(req dto.BookingRequest) (*dto.BookingResponse, error) {
	booking := u.repo.Create(req)
	u.cache.Set(booking.ID, booking)
	return booking, nil
}

// Get booking by id
func (u *bookingUsecase) GetBookingByID(id int) (*dto.BookingResponse, error) {
	// try get data from cache
	cachedBooking, err := u.cache.Get(id)
	if err == nil {
		return cachedBooking, nil
	}

	// get data from repository
	booking, exists := u.repo.GetByID(id)
	if !exists {
		return nil, fmt.Errorf("booking not found")
	}

	// set data to cache
	u.cache.Set(id, booking)

	return booking, nil
}

// Get all bookings
func (u *bookingUsecase) GetAllBookings(sortParam, highValue string) ([]*dto.BookingResponse, error) {
	var bookings []*dto.BookingResponse

	if highValue == "true" {
		bookings = u.repo.GetHighValueBookings(50000)
	} else {
		bookings = u.repo.GetAll()
	}

	// เรียงลำดับ
	switch sortParam {
	case "price":
		sort.Slice(bookings, func(i, j int) bool {
			return bookings[i].Price < bookings[j].Price
		})
	case "date":
		sort.Slice(bookings, func(i, j int) bool {
			t1, _ := time.Parse(time.RFC3339, bookings[i].CreatedAt)
			t2, _ := time.Parse(time.RFC3339, bookings[j].CreatedAt)
			return t1.Before(t2)
		})
	}

	return bookings, nil
}

// Update status booking
func (u *bookingUsecase) UpdateBooking(id int, status string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.repo.UpdateBookingStatus(id, status)
	if err != nil {
		return fmt.Errorf("booking not found")
	}

	// get data from repository
	booking, exists := u.repo.GetByID(id)
	if !exists {
		return fmt.Errorf("booking not found in repository")
	}

	// update cache
	u.cache.Set(id, booking)

	return nil
}

// Cancel booking
func (u *bookingUsecase) CancelBooking(id int) error {
	// change status to canceled
	err := u.repo.UpdateBookingStatus(id, "canceled")
	if err != nil {
		return fmt.Errorf("failed to cancel booking")
	}

	// delete from cache
	u.cache.Delete(id)

	return nil
}

// Background task for check expired booking
func (u *bookingUsecase) BackgroundTaskBooking(wg *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Minute)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range ticker.C {
			u.checkExpiredBookings()
		}
	}()
}

// checkExpiredBookings handles the expiration logic for pending bookings
func (u *bookingUsecase) checkExpiredBookings() {
	bookings := u.repo.GetAll()
	currentTime := time.Now()

	for _, booking := range bookings {
		if booking.Status == "pending" {
			createdAt, err := time.Parse(time.RFC3339, booking.CreatedAt)
			if err != nil {
				continue // ข้ามถ้ามี error ในการ parse เวลา
			}
			if currentTime.Sub(createdAt) > 5*time.Minute {
				// เปลี่ยนสถานะเป็น canceled
				err := u.repo.UpdateBookingStatus(booking.ID, "canceled")
				if err != nil {
					continue // ข้ามถ้าไม่สามารถอัปเดต
				}
				// อัปเดตแคช
				booking.Status = "canceled"
				u.cache.Set(booking.ID, booking)
			}
		}
	}
}

// Update booking status
func (u *bookingUsecase) UpdateBookingStatus(id int, status string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Change status in Repository
	err := u.repo.UpdateBookingStatus(id, status)
	if err != nil {
		return fmt.Errorf("failed to update booking status")
	}

	// Update cache if it exists
	cachedBooking, err := u.cache.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get booking from cache")
	}

	cachedBooking.Status = status
	cachedBooking.UpdatedAt = time.Now().Format(time.RFC3339)
	u.cache.Set(id, cachedBooking)

	return nil
}

func (u *bookingUsecase) CheckExpiredBookings() {
    bookings := u.repo.GetAll()
    for _, booking := range bookings {
        if booking.Status == "pending" && isExpired(booking.CreatedAt) {
            u.repo.UpdateBookingStatus(booking.ID, "canceled")
            u.cache.Delete(booking.ID)
        }
    }
}


func isExpired(createdAt string) bool {
    return false
}
