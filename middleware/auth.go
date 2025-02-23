package middleware

import (
	"os"
	"strings"

	"github.com/Eursukkul/fiber-booking-system/utils"
	"github.com/gofiber/fiber/v2"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middlware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middlware-002"
	apiKeyErr      middlewareHandlersErrCode = "middlware-005"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")

		if token == "" {
			return utils.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(jwtAuthErr),
				"Missing or invalid token",
			).Res()
		}
		claims, err := utils.ParseToken(os.Getenv("JWT_SECRET"), token)
		if err != nil {
			return utils.NewResponse(c).Error(
				fiber.StatusUnauthorized,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		c.Locals("claims", claims)

		return c.Next()
	}
}

func (m *AuthMiddleware) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Api-Key")
		if _, err := utils.ParseApiKey(os.Getenv("API_KEY"), key); err != nil {
			return utils.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(apiKeyErr),
				"apikey is invalid or required",
			).Res()
		}
		return c.Next()
	}
}