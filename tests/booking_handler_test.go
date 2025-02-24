package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Eursukkul/fiber-booking-system/dto"
	"github.com/Eursukkul/fiber-booking-system/handler"
	"github.com/Eursukkul/fiber-booking-system/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp(handler *handler.BookingHandler) *fiber.App {
	app := fiber.New()
	app.Post("/api/bookings", handler.CreateBooking)
	app.Get("/api/bookings/:id", handler.GetBookingByID)
	app.Get("/api/bookings", handler.GetAllBookings)
	app.Delete("/api/bookings/:id", handler.CancelBooking)
	return app
}

func TestCreateBooking_Success(t *testing.T) {
	mockUsecase := new(mocks.MockBookingUsecase)
	bookingHandler := handler.NewBookingHandler(mockUsecase)
	app := setupTestApp(bookingHandler)

	reqBody := dto.BookingRequest{
		UserID:    1,
		ServiceID: 2,
		Price:     1000,
	}

	expectedResp := &dto.BookingResponse{
		ID:        1,
		UserID:    1,
		ServiceID: 2,
		Price:     1000,
		Status:    "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	mockUsecase.On("CreateBooking", reqBody).Return(expectedResp, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/bookings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dto.BookingResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, expectedResp.ID, response.ID)
	assert.Equal(t, expectedResp.Status, response.Status)

	mockUsecase.AssertExpectations(t)
}

func TestGetBookingByID_Success(t *testing.T) {
	mockUsecase := new(mocks.MockBookingUsecase)
	bookingHandler := handler.NewBookingHandler(mockUsecase)
	app := setupTestApp(bookingHandler)

	expectedResp := &dto.BookingResponse{
		ID:        1,
		UserID:    1,
		ServiceID: 2,
		Price:     1000,
		Status:    "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	mockUsecase.On("GetBookingByID", 1).Return(expectedResp, nil)

	req := httptest.NewRequest("GET", "/api/bookings/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dto.BookingResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, expectedResp.ID, response.ID)

	mockUsecase.AssertExpectations(t)
}

func TestGetBookingByID_NotFound(t *testing.T) {
	mockUsecase := new(mocks.MockBookingUsecase)
	bookingHandler := handler.NewBookingHandler(mockUsecase)
	app := setupTestApp(bookingHandler)

	mockUsecase.On("GetBookingByID", 999).Return(nil, errors.New("booking not found"))

	req := httptest.NewRequest("GET", "/api/bookings/999", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}

func TestGetAllBookings_Success(t *testing.T) {
	mockUsecase := new(mocks.MockBookingUsecase)
	bookingHandler := handler.NewBookingHandler(mockUsecase)
	app := setupTestApp(bookingHandler)

	expectedResp := []*dto.BookingResponse{
		{
			ID:        1,
			UserID:    1,
			ServiceID: 2,
			Price:     1000,
			Status:    "pending",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        2,
			UserID:    2,
			ServiceID: 3,
			Price:     2000,
			Status:    "confirmed",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}

	mockUsecase.On("GetAllBookings", "price", "false").Return(expectedResp, nil)

	req := httptest.NewRequest("GET", "/api/bookings?sort=price&high-value=false", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response []*dto.BookingResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, len(expectedResp), len(response))

	mockUsecase.AssertExpectations(t)
}

func TestCancelBooking_Success(t *testing.T) {
	mockUsecase := new(mocks.MockBookingUsecase)
	bookingHandler := handler.NewBookingHandler(mockUsecase)
	app := setupTestApp(bookingHandler)

	mockBooking := &dto.BookingResponse{
		ID:     1,
		Status: "pending",
	}

	mockUsecase.On("GetBookingByID", 1).Return(mockBooking, nil)
	mockUsecase.On("CancelBooking", 1).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/bookings/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}

func TestCancelBooking_AlreadyConfirmed(t *testing.T) {
	mockUsecase := new(mocks.MockBookingUsecase)
	bookingHandler := handler.NewBookingHandler(mockUsecase)
	app := setupTestApp(bookingHandler)

	mockBooking := &dto.BookingResponse{
		ID:     1,
		Status: "confirmed",
	}

	mockUsecase.On("GetBookingByID", 1).Return(mockBooking, nil)

	req := httptest.NewRequest("DELETE", "/api/bookings/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}
