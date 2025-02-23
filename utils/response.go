package utils

import (
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, traceId, msg string) IResponse
	Res() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    *fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	Msg string `json:"message"`
}

func NewResponse(c *fiber.Ctx) IResponse {
	return &Response{
		Context: c,
	}
}

func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	return r
}

func (r *Response) Error(code int, traceId, msg string) IResponse {
	_ = traceId
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		Msg: msg,
	}
	r.IsError = true
	return r
}
func (r *Response) Res() error {
	if r.IsError {
		return r.Context.Status(r.StatusCode).JSON(r.ErrorRes)
	}
	return r.Context.Status(r.StatusCode).JSON(r.Data)
}
