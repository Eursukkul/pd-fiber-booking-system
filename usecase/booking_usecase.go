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
	BookingUsecase interface{}

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
func (u *bookingUsecase) Create(req dto.BookingRequest) (*dto.BookingResponse, error) {
	booking := u.repo.Create(req)
	u.cache.Set(booking.ID, booking)
	return booking, nil
}

// Get booking by id
func (u *bookingUsecase) GetByID(id int) (*dto.BookingResponse, error) {
	return u.cache.Get(id)
}

// Get all bookings
func (u *bookingUsecase) GetAll(sortBy string, highValue bool) ([]*dto.BookingResponse, error) {
	bookings := u.repo.GetAll()

	var filteredBookings []*dto.BookingResponse
	for _, booking := range bookings {
		if highValue && booking.Price > 50000 {
			continue
		}
		filteredBookings = append(filteredBookings, booking)
	}

	switch sortBy {
	case "price":
		sort.Slice(filteredBookings, func(i, j int) bool {
			return filteredBookings[i].Price < filteredBookings[j].Price
		})
	case "date":
		sort.Slice(filteredBookings, func(i, j int) bool {
			timeI, err := time.Parse(time.RFC3339, filteredBookings[i].CreatedAt)
			if err != nil {
				return false
			}
			timeJ, err := time.Parse(time.RFC3339, filteredBookings[j].CreatedAt)
			if err != nil {
				return false
			}
			return timeI.Before(timeJ)
		})
	default:
		sort.Slice(filteredBookings, func(i, j int) bool {
			return filteredBookings[i].ID < filteredBookings[j].ID
		})
	}
	return filteredBookings, nil
}

// Update status booking
func (u *bookingUsecase) Update(id int, status string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	updated := u.repo.Update(id, status)
	if !updated {
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
func (u *bookingUsecase) Cancel(id int) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	booking, err := u.repo.GetByID(id)
	if !err {
		return fmt.Errorf("booking not found")
	}

	if booking.Status == "confirmed" {
		return fmt.Errorf("cannot cancel confirmed booking")
	}

	updated := u.repo.Update(id, "canceled")
	if !updated {
		return fmt.Errorf("failed to cancel booking")
	}

	u.cache.Delete(id)
	return nil
}

// Background task for check expired booking
func (u *bookingUsecase) BackgroundTask(wg *sync.WaitGroup) {
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
	for _, booking := range bookings {
		if booking.Status != "pending" {
			continue
		}

		createdAt, err := time.Parse(time.RFC3339, booking.CreatedAt)
		if err != nil {
			continue
		}

		if time.Since(createdAt) > 5*time.Minute {
			u.UpdateBookingStatus(booking.ID, "canceled")
		}
	}
}

// Update booking status
func (u *bookingUsecase) UpdateBookingStatus(id int, status string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Change status in Repository
	updated := u.repo.Update(id, status)
	if !updated {
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
