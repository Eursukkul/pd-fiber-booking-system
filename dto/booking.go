package dto

type (
	BookingRequest struct {
		UserID    int     `json:"user_id" validate:"required"`
		ServiceID int     `json:"service_id" validate:"required"`
		Price     float64 `json:"price" validate:"required,gt=0"`
	}

	BookingResponse struct {
		ID        int     `json:"id"`
		UserID    int     `json:"user_id"`
		ServiceID int     `json:"service_id"`
		Price     float64 `json:"price"`
		Status    string  `json:"status"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}

	// SwaggerResponse represents a standard API response
	SwaggerResponse struct {
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}

	// ErrorResponse represents an error response
	ErrorResponse struct {
		Message string `json:"message"`
	}
)
