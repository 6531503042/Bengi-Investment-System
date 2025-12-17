package middleware

import (
	"strings"

	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string
		// 1. ลองดึงจาก Cookie ก่อน (สำหรับ Web)
		token = c.Cookies(utils.AccessTokenCookie)
		// 2. ถ้าไม่มี Cookie ให้ดึงจาก Authorization header (สำหรับ Mobile/API)
		if token == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}
		}
		// 3. ถ้าไม่มี token เลย
		if token == "" {
			return common.Unauthorized(c, "Authentication required")
		}
		// 4. Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			// ถ้า token expired, ลอง clear cookie
			utils.ClearAuthCookies(c)
			return common.Unauthorized(c, "Invalid or expired token")
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
