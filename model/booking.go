package models

import "time"

type (
	Booking struct {
		ID        int
		UserID    int
		ServiceID int
		Price     float64
		Status    string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
