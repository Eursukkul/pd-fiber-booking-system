package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoggerMiddleware struct {
}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (m *LoggerMiddleware) Logger(c *fiber.Ctx) error {
	
	startTime := time.Now()
	err := c.Next()
	stop := time.Now()

	log.Printf("Method: %s, Path: %s, Status: %d, Duration: %s", c.Method(), c.Path(), c.Response().StatusCode(), stop.Sub(startTime))

	return err
}
