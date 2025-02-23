package handler

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/Eursukkul/fiber-booking-system/dto"
	"github.com/Eursukkul/fiber-booking-system/usecase"
	"github.com/gofiber/fiber/v2"
)

type (
	BookingHandler struct {
		BookingUsecase usecase.BookingUsecase
	}
)

func NewBookingHandler(BookingUsecase usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{BookingUsecase: BookingUsecase}
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking with user_id, service_id, and price
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body dto.BookingRequest true "Booking Request"
// @Success 201 {object} dto.SwaggerResponse{data=dto.BookingResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var req dto.BookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	booking, err := h.BookingUsecase.CreateBooking(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if booking.Price > 50000 {
		go func(id int) {
			//Random status confirm or rejected
			status := "confirmed"
			if rand.Intn(2) == 0 {
				status = "rejected"
			}
			time.Sleep(time.Second * 1) //Delay 1 second
			//Update status
			err := h.BookingUsecase.UpdateBookingStatus(id, status)
			if err != nil {
				log.Printf("Failed to update booking status: %v", err)
			}
		}(booking.ID)
	}

	return c.Status(fiber.StatusCreated).JSON(booking)
}

// GetBookingByID godoc
// @Summary Get a booking by ID
// @Description Get booking details by ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} dto.SwaggerResponse{data=dto.BookingResponse}
// @Failure 404 {object} dto.ErrorResponse
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBookingByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	booking, err := h.BookingUsecase.GetBookingByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Booking not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(booking)
}

// GetAllBookings godoc
// @Summary Get all bookings
// @Description Get all bookings with optional sorting and filtering
// @Tags bookings
// @Accept json
// @Produce json
// @Param sort query string false "Sort by (price or date)"
// @Param high-value query bool false "Filter high value bookings"
// @Success 200 {object} dto.SwaggerResponse{data=[]dto.BookingResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /bookings [get]
func (h *BookingHandler) GetAllBookings(c *fiber.Ctx) error {
	sortBy := c.Query("sort", "id")
	highValue := c.Query("high-value")
	bookings, err := h.BookingUsecase.GetAllBookings(sortBy, highValue)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve bookings",
		})
	}

	return c.JSON(bookings)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel an existing booking by ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} dto.SwaggerResponse
// @Failure 400,404 {object} dto.ErrorResponse
// @Router /bookings/{id} [delete]
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid booking ID",
		})
	}

	// ตรวจสอบสถานะก่อนยกเลิก
	booking, err := h.BookingUsecase.GetBookingByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Booking not found",
		})
	}

	if booking.Status == "confirmed" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot cancel a confirmed booking",
		})
	}

	// ยกเลิกการจอง
	err = h.BookingUsecase.CancelBooking(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Booking canceled successfully",
	})
}
