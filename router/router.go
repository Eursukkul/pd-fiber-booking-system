package router

import (
	"github.com/Eursukkul/fiber-booking-system/handler"
	"github.com/Eursukkul/fiber-booking-system/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, bookingHandler *handler.BookingHandler, logger *middleware.LoggerMiddleware) {
	api := app.Group("/api")

	api.Post("/bookings", logger.Logger, bookingHandler.CreateBooking)
	api.Get("/bookings/:id", logger.Logger, bookingHandler.GetBookingByID)
	api.Get("/bookings", logger.Logger, bookingHandler.GetAllBookings)
	api.Delete("/bookings/:id", logger.Logger, bookingHandler.CancelBooking)
}
// if use middleware auth
func SetupRoutes_middleware(app *fiber.App, bookingHandler *handler.BookingHandler, logger *middleware.LoggerMiddleware, auth *middleware.AuthMiddleware) {
	api := app.Group("/v1")
 	
	api.Post("/bookings", auth.JwtAuth(), logger.Logger, bookingHandler.CreateBooking)
	api.Get("/bookings/:id", auth.JwtAuth(), logger.Logger, bookingHandler.GetBookingByID)
	api.Get("/bookings", auth.JwtAuth(), logger.Logger, bookingHandler.GetAllBookings)
	api.Delete("/bookings/:id", auth.JwtAuth(), logger.Logger, bookingHandler.CancelBooking)
}