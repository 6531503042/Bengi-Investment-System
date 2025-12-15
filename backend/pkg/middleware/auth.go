package middleware

import (
	"strings"

	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return common.Unauthorized(c, "Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return common.Unauthorized(c, "Invalid authorization header")
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			return common.Unauthorized(c, "Invalid authorization header")
		}

		c.Locals("userId", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("roleId", claims.RoleID)

		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) string {
	if id, ok := c.Locals("userId").(string); ok {
		return id
	}
	return ""
}
